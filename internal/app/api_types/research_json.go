package apitypes

import (
	"LAB1/internal/app/ds"
	"time"
)

type ResearchJSON struct {
	ID           	int       		`json:"id"`
	Status       	string    		`json:"status"`
	DateResearch 	string 		`json:"date_research"`               
	DateCreate   	time.Time 		`json:"date_create"`                 
	CreatorLogin    string       	`json:"creator_login"`           
	DateForm     	*time.Time 		`json:"date_form"`              
	DateFinish   	*time.Time 		`json:"date_finish"`                
	ModeratorLogin  *string       	`json:"moderator_login"`                 

}

func ResearchToJSON(research ds.Research, creatorLogin string, moderatorLogin string) ResearchJSON {
	var dateForm, dateFinish *time.Time
	if research.DateForm.Valid {
		dateForm = &research.DateForm.Time
	}

	if research.DateFinish.Valid {
		dateFinish = &research.DateFinish.Time
	}

	var mLogin *string
	if moderatorLogin != "" {
		mLogin = &moderatorLogin
	}





	return ResearchJSON{
		ID:				research.ID,
		Status:       	research.Status,
		DateResearch: 	research.DateResearch,              
		DateCreate:   	research.DateCreate,                 
		CreatorLogin:   creatorLogin,           
		DateForm:     	dateForm,              
		DateFinish:   	dateFinish,               
		ModeratorLogin:	mLogin,
	}
}


func ResearchFromJSON(research ResearchJSON) ds.Research {
	if research.DateResearch == "" {
		return ds.Research{}
	}
	return ds.Research{
		DateResearch: research.DateResearch,
	}
}

type StatusJSON struct {
	Status string `json:"status"`
}