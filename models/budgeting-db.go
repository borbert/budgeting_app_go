package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

//Returns one budget item and error if any
func (m *DBModel) Get(id int) (*BudgetAcct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select *
	from budget_acct ba
	where ba.id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var ba BudgetAcct

	err := row.Scan(
		&ba.ID,
		&ba.UserID,
		&ba.Item,
		&ba.Description,
		&ba.Amount,
		&ba.BudgetingType,
		&ba.Biweekly,
		&ba.ApplyDefault,
		&ba.DefaultAmt,
		&ba.CreatedAt,
		&ba.UpdatedAt,
		&ba.TerminatedAt,
	)

	if err != nil {
		return nil, err
	}

	//get tags if any
	query = `select  *
	from budget_tags bt 
	where bt.item_id = $1
	order by bt.tag_id`

	rows, _ := m.DB.QueryContext(ctx, query, id)
	defer rows.Close()

	tags := make(map[int]string)
	for rows.Next() {
		var bt BudgetTag
		err := rows.Scan(
			&bt.ID,
			&bt.TagID,
			&bt.ItemID,
			&bt.Description,
		)
		if err != nil {
			return nil, err
		}
		tags[bt.TagID] = bt.Description
	}

	ba.Tags = tags

	return &ba, nil
}

//Returns all movies and error if any
func (m *DBModel) All(tag ...int) ([]*BudgetAcct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// fmt.Println(tag)

	where := ""
	if len(tag) > 0 {
		where = fmt.Sprintf("where id in (select item_id from budget_tags where tag_id = %d)", tag[0])
	}

	query := fmt.Sprintf(`select * from budget_acct %s order by id`, where)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budget_accts []*BudgetAcct

	for rows.Next() {
		var budget_acct BudgetAcct
		err := rows.Scan(
			&budget_acct.ID,
			&budget_acct.UserID,
			&budget_acct.Item,
			&budget_acct.Description,
			&budget_acct.Amount,
			&budget_acct.BudgetingType,
			&budget_acct.Biweekly,
			&budget_acct.ApplyDefault,
			&budget_acct.DefaultAmt,
			&budget_acct.CreatedAt,
			&budget_acct.UpdatedAt,
			&budget_acct.TerminatedAt,
		)
		if err != nil {
			return nil, err
		}
		//get tags if any
		tagQuery := `select  *
		from budget_tags bt 
		where bt.item_id=$1`

		tagRows, _ := m.DB.QueryContext(ctx, tagQuery, budget_acct.ID)
		defer tagRows.Close()

		tags := make(map[int]string)
		for tagRows.Next() {
			var bt BudgetTag
			err := tagRows.Scan(
				&bt.ID,
				&bt.TagID,
				&bt.ItemID,
				&bt.Description,
			)
			if err != nil {
				return nil, err
			}
			tags[bt.TagID] = bt.Description
		}

		budget_acct.Tags = tags

		budget_accts = append(budget_accts, &budget_acct)
	}

	return budget_accts, nil
}

func (m *DBModel) AllTags() ([]*Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select * from tags order by description`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*Tag

	for rows.Next() {
		var t Tag
		err := rows.Scan(
			&t.TagID,
			&t.Description,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &t)
	}
	return tags, nil
}

func (m *DBModel) InsertBudgetItem(budgetacct BudgetAcct) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into budget_acct 
	(id, user_id, item, description, amount, budgeting_type, biweekly,
		apply_default_amount, default_amount, created_at, updated_at, terminated_at) 
	values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		budgetacct.ID,
		budgetacct.UserID,
		budgetacct.Item,
		budgetacct.Description,
		budgetacct.Amount,
		budgetacct.BudgetingType,
		budgetacct.Biweekly,
		budgetacct.ApplyDefault,
		budgetacct.DefaultAmt,
		budgetacct.CreatedAt,
		budgetacct.UpdatedAt,
		budgetacct.TerminatedAt,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *DBModel) UpdateBudgetItem(budgetacct BudgetAcct) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update budget_acct set
				item=$1, 
				description=$2, 
				amount=$3, 
				budgeting_type=$4, 
				biweekly=$5,
				apply_default_amount=$6, 
				default_amount=$7, 
				updated_at=$8, 
			where id=$9 and user_id=$10
			`

	_, err := m.DB.ExecContext(ctx, stmt,
		budgetacct.Item,
		budgetacct.Description,
		budgetacct.Amount,
		budgetacct.BudgetingType,
		budgetacct.Biweekly,
		budgetacct.ApplyDefault,
		budgetacct.DefaultAmt,
		budgetacct.UpdatedAt,
		budgetacct.ID,
		budgetacct.UserID,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *DBModel) DeleteBudgetItem(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "delete from budget_acct where id=$1"

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *DBModel) AllUsers() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "select * from users"

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// get one user
func (m *DBModel) GetOneUser(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select *
	from users u
	where u.email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)

	var u User

	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Password,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
