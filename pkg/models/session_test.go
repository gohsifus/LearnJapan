package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSession_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil{
		panic(err)
	}
	defer db.Close()

	DB = db

	type args struct{
		session Session
	}

	type mockBehaviour func(args args, expectResult interface{})

	testTable := []struct{
		name string
		mockBehaviour mockBehaviour
		args args
		expectResult interface{}
		wantErr bool
	} {
		{
			name: "OK if SessionId not exist",
			args: args{
				session: Session{
					SessionId: "qwel",
					UserId: 13,
					Expires: "2022-05-03 12:10:20",
				},
			},
			expectResult: true,
			mockBehaviour: func(args args, expectResult interface{}){

				rows := sqlmock.NewRows([]string{""}).AddRow(expectResult)
				mock.ExpectQuery("INSERT INTO sessions").
					WithArgs(args.session.SessionId, args.session.UserId, args.session.Expires).WillReturnRows(rows)
			},
		},
	}

	for _, testCase := range testTable{
		t.Run(testCase.name, func(t *testing.T){
			testCase.mockBehaviour(testCase.args, testCase.expectResult)

			got, err := testCase.args.session.Add()
			if testCase.wantErr{ //Если ожидаем ошибку
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectResult, got)
			}
		})
	}
}