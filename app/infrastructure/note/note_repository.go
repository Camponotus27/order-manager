package note

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"order-manager/app/domain/model"
	"order-manager/app/shared/utils"

	"github.com/dstotijn/go-notion"
)

const (
	propertyRemaining = "Restante"
	propertySold      = "Vendidos"
	propertyDate      = "Fecha"
)

const maxLengthFileNameDefault = 100
const maxLengthFileNameCaption = 100

const propertyID = "ID"
const propertyLastEditedTime = "Last edited time"

type DBConfig struct {
	IDDBOrder string
}

type Repository struct {
	NotionClient *notion.Client
	dbConfig     *DBConfig

	filterHaveRemaining bool

	orderDateDesc bool
}

func NewListItemRepository(token string, dbConfig *DBConfig, client *http.Client) *Repository {
	notionClient := notion.NewClient(token, notion.WithHTTPClient(client))
	return &Repository{
		NotionClient: notionClient,
		dbConfig:     dbConfig,
	}
}

func (c *Repository) Current(ctx context.Context) (int, error) {
	cookies, err := c.orderRemainingCookies(ctx)
	if err != nil {
		return 0, err
	}
	return cookieRemainingTotal(cookies), nil
}

func (c *Repository) SellOne(ctx context.Context) (int, error) {
	cookie, errGetCookie := c.OrderDateDesc().HasRemaining().First(ctx)
	if errGetCookie != nil {
		return 0, errGetCookie
	}

	IdDB := c.dbConfig.IDDBOrder
	if err := c.validateDBandClient(IdDB); err != nil {
		return 0, err
	}

	remainingMinusOne := float64(cookie.Sold + 1)
	paramUpdate := notion.UpdatePageParams{
		DatabasePageProperties: map[string]notion.DatabasePageProperty{
			propertySold: {
				Number: &remainingMinusOne,
			},
		},
	}
	_, errUpdate := c.NotionClient.UpdatePage(ctx, cookie.UID, paramUpdate)
	if errUpdate != nil {
		return 0, errUpdate
	}
	return c.Current(ctx)
}

func (c *Repository) orderRemainingCookies(ctx context.Context) ([]*model.CookieItem, error) {
	return c.HasRemaining().Get(ctx)
}

func (c *Repository) First(ctx context.Context) (*model.CookieItem, error) {
	cookies, err := c.Get(ctx)
	if err != nil {
		return nil, err
	}

	if len(cookies) == 0 {
		return nil, errors.New("no found cookie")
	}

	return cookies[0], nil
}

func (c *Repository) OrderDateDesc() *Repository {
	c.orderDateDesc = true
	return c
}

func (c *Repository) HasRemaining() *Repository {
	c.filterHaveRemaining = true
	return c
}

func (c *Repository) validateDBandClient(IdDB string) error {
	if IdDB == "" {
		return errors.New("ID DB is no provided")
	}

	if c.NotionClient == nil {
		return errors.New("notion client is no provided")
	}
	return nil
}

func (c *Repository) Get(ctx context.Context) ([]*model.CookieItem, error) {
	zero := 0

	IdDB := c.dbConfig.IDDBOrder
	if err := c.validateDBandClient(IdDB); err != nil {
		return nil, err
	}

	var filters []notion.DatabaseQueryFilter
	if c.filterHaveRemaining == true {
		filters = append(filters, notion.DatabaseQueryFilter{
			Property: propertyRemaining,
			DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
				Number: &notion.NumberDatabaseQueryFilter{
					GreaterThan: &zero,
				},
			},
		})
	}

	var order []notion.DatabaseQuerySort
	if c.orderDateDesc == true {
		order = append(order, notion.DatabaseQuerySort{
			Property:  propertyDate,
			Direction: notion.SortDirDesc,
		})
	}

	query := &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			And: filters,
		},
		Sorts: order,
	}
	result, err := c.NotionClient.QueryDatabase(ctx, IdDB, query)

	if err != nil {
		return nil, err
	}

	return mapPagesToListItems(result.Results), nil
}

func (c *Repository) getObject(ctx context.Context, objectNotion *model.ObjectNotion) (*model.ObjectNotion, error) {
	if c.NotionClient == nil {
		return nil, errors.New("notion client is no provided")
	}

	if objectNotion == nil || objectNotion.ID == "" {
		return nil, errors.New("ID of object is not provided")
	}

	switch objectNotion.Type {
	case model.PageID:
		result, err := c.NotionClient.FindPageByID(ctx, objectNotion.ID)
		if err != nil {
			return nil, err
		}
		return mapPageToObjectNotion(result), nil
	case model.DatabaseID:
		result, err := c.NotionClient.FindDatabaseByID(ctx, objectNotion.ID)
		if err != nil {
			return nil, err
		}

		return mapDataBaseToObjectNotion(result), nil
	case model.BlockID:
		return nil, errors.New("object notion type not supported for the momentum")
	default:
		return nil, errors.New("object notion type not supported")
	}
}

func (c *Repository) findObjectByID(ctx context.Context, task *model.Note) (notion.Page, error) {
	return c.NotionClient.FindPageByID(ctx, task.ID)
}

func (c *Repository) getObjectWithBody(ctx context.Context, objectNotion *model.ObjectNotion) (*model.ObjectNotion, error) {
	objectNotionResult, errGetObj := c.getObject(ctx, objectNotion)
	if errGetObj != nil {
		return nil, errGetObj
	}

	files, errGetBlocks := c.getBlockChildren(ctx, objectNotionResult)
	if errGetBlocks != nil {
		return nil, errGetBlocks
	}
	objectNotionResult.Files = files
	return objectNotionResult, errGetObj

}

func (c *Repository) getBlockChildren(ctx context.Context, objectNotion *model.ObjectNotion) ([]*model.File, error) {
	result, err := c.NotionClient.FindBlockChildrenByID(ctx, objectNotion.ID, nil)
	if err != nil {
		return nil, err
	}

	return getFile(objectNotion.Name, result.Results)

}

func getFile(fileNameDefault string, blocks []notion.Block) ([]*model.File, error) {
	var files []*model.File
	for _, block := range blocks {
		switch t := block.(type) {
		case *notion.ImageBlock:
			if t.File == nil || t.File.URL == "" {
				return nil, errors.New("URL is empty")
			}

			nameFile := utils.GetFirstN(maxLengthFileNameDefault, fileNameDefault)
			strCaption := toStr(t.Caption)
			if strCaption != "" {
				nameFile = fmt.Sprintf("%s - %s", nameFile, utils.GetFirstN(maxLengthFileNameCaption, strCaption))
			}

			files = append(files, &model.File{
				Url:  t.File.URL,
				Name: nameFile,
			})
		case *notion.BulletedListItemBlock:
			f, err := getFile(fileNameDefault, t.Children)
			if err != nil {
				return nil, err
			}
			files = append(files, f...)
		case *notion.ParagraphBlock:
			f, err := getFile(fileNameDefault, t.Children)
			if err != nil {
				return nil, err
			}
			files = append(files, f...)
		case *notion.NumberedListItemBlock:
			f, err := getFile(fileNameDefault, t.Children)
			if err != nil {
				return nil, err
			}
			files = append(files, f...)
		case *notion.QuoteBlock:
			f, err := getFile(fileNameDefault, t.Children)
			if err != nil {
				return nil, err
			}
			files = append(files, f...)
		case *notion.ToggleBlock:
			f, err := getFile(fileNameDefault, t.Children)
			if err != nil {
				return nil, err
			}
			files = append(files, f...)
		case *notion.TemplateBlock:
			f, err := getFile(fileNameDefault, t.Children)
			if err != nil {
				return nil, err
			}
			files = append(files, f...)
		case *notion.Heading1Block:
			f, err := getFile(fileNameDefault, t.Children)
			if err != nil {
				return nil, err
			}
			files = append(files, f...)
		case *notion.Heading2Block:
			f, err := getFile(fileNameDefault, t.Children)
			if err != nil {
				return nil, err
			}
			files = append(files, f...)
		case *notion.Heading3Block:
			f, err := getFile(fileNameDefault, t.Children)
			if err != nil {
				return nil, err
			}
			files = append(files, f...)
		case *notion.ToDoBlock:
			f, err := getFile(fileNameDefault, t.Children)
			if err != nil {
				return nil, err
			}
			files = append(files, f...)
		case *notion.CalloutBlock:
			f, err := getFile(fileNameDefault, t.Children)
			if err != nil {
				return nil, err
			}
			files = append(files, f...)
		case *notion.CodeBlock:
			f, err := getFile(fileNameDefault, t.Children)
			if err != nil {
				return nil, err
			}
			files = append(files, f...)
		case *notion.ColumnBlock:
			f, err := getFile(fileNameDefault, t.Children)
			if err != nil {
				return nil, err
			}
			files = append(files, f...)
		case *notion.SyncedBlock:
			f, err := getFile(fileNameDefault, t.Children)
			if err != nil {
				return nil, err
			}
			files = append(files, f...)
		}
	}

	return files, nil
}
