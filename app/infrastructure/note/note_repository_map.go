package note

import (
	"fmt"
	"time"

	"order-manager/app/domain/model"

	"github.com/dstotijn/go-notion"
)

func mapPagesToListItems(pages []notion.Page) []*model.CookieItem {
	var components []*model.CookieItem

	for _, page := range pages {
		components = append(components, mapPageToListItem(page))
	}

	return components
}

func mapPageToListItem(page notion.Page) *model.CookieItem {
	properties, ok := page.Properties.(notion.DatabasePageProperties)
	if !ok {
		return &model.CookieItem{
			UID:   time.Now().String(),
			Title: "Properties is not DatabasePageProperties",
		}
	}

	valueStr := valueStrPageProperty(properties["Name"])
	if valueStr == "" {
		valueStr = valueStrPageProperty(properties["Title"])
	}

	remaining := valueIntPageProperty(properties["Restante"])
	sold := valueIntPageProperty(properties["Vendidos"])

	return &model.CookieItem{
		UID:       page.ID,
		Title:     valueStr,
		Subtitle:  valueStr,
		Arg:       page.ID,
		Remaining: int(remaining),
		Sold:      int(sold),
	}
}
func mapPageToObjectNotion(page notion.Page) *model.ObjectNotion {
	return model.NewObjNotFromAPI().
		SetID(page.ID).
		SetName(nameFromPage(page)).
		SetObjNotParent(getObjNotFromParent(page.Parent))
}

func cookieRemainingTotal(cookies []*model.CookieItem) int {
	total := 0
	for _, cookie := range cookies {
		total += cookie.Remaining
	}
	return total
}

func nameFromPage(page notion.Page) string {
	switch p := page.Properties.(type) {
	case notion.DatabasePageProperties:
		valueStr := valueStrPageProperty(p["Name"])
		if valueStr != "" {
			return valueStr
		}
		return valueStrPageProperty(p["Title"])
	case notion.PageProperties:
		return toStr(p.Title.Title)
	default:
		return "Not is possible get name from page"
	}
}

func mapDataBaseToObjectNotion(database notion.Database) *model.ObjectNotion {
	return model.NewObjNotFromAPI().
		SetID(database.ID).
		SetName(toStr(database.Title)).
		SetObjNotParent(getObjNotFromParent(database.Parent))
}

func getObjNotFromParent(parent notion.Parent) *model.ObjectNotion {
	switch string(parent.Type) {
	case string(model.PageID):
		return &model.ObjectNotion{
			ID:   parent.PageID,
			Type: model.PageID,
		}
	case string(model.BlockID):
		return &model.ObjectNotion{
			ID:   parent.BlockID,
			Type: model.BlockID,
		}
	case string(model.DatabaseID):
		return &model.ObjectNotion{
			ID:   parent.DatabaseID,
			Type: model.DatabaseID,
		}
	default:
		return nil
	}
}

func valueIntPageProperty(prop notion.DatabasePageProperty) float64 {
	value := prop.Value()

	switch t := value.(type) {
	case *notion.FormulaResult:
		if t.Number == nil {
			return 0
		}
		return *t.Number
	case *float64:
		return *t
	default:
		return 0
	}
}

func valueStrPageProperty(prop notion.DatabasePageProperty) string {
	value := prop.Value()

	switch t := value.(type) {
	case []notion.RichText:
		return toStr(t)
	default:
		return ""
	}
}

func toStr(richTexts []notion.RichText) string {
	var text string
	for _, richText := range richTexts {
		if text == "" {
			text = richText.PlainText
			continue
		}

		text = fmt.Sprintf("%s %s", text, richText.PlainText)
	}
	return text
}
