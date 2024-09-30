package user

type User struct {
	Id        int    `db:"id" json:"id"`
	Address   string `db:"address" json:"address"`
	InviterId int    `db:"inviter_id" json:"inviter_id"`
	Points    int    `db:"points" json:"points"`
}
