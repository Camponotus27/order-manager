package model

type MenuOption string

const (
	UpdateArtifacts                 MenuOption = "update_artifacts"
	ShowArtifacts                   MenuOption = "show_artifacts"
	MenuContext                     MenuOption = "menu_context"
	SaveContext                     MenuOption = "save_context"
	RenameFileWithContext           MenuOption = "rename_file_with_context"
	SearchTask                      MenuOption = "search_task"
	DownloadNote                    MenuOption = "download_note"
	ParseResponseToPostNotification MenuOption = "parse_response_to_post_notification"
	Menu                            MenuOption = "menu"
)
