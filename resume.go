package parser

type Resume struct {
	Id             int      `json:"id"`
	First_name     string   `json:"first_name"`
	Middle_name    string   `json:"middle_name"`
	Last_name      string   `json:"last_name"`
	Specialization []string `json:"specialization"`
	Age            int      `json:"age"`
	Experience     int      `json:"experience"`
}
