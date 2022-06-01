package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSession_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil{
		t.Error("Ошибка создания mock")
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
			name: "OK",
			args: args{
				session: Session{
					SessionId: "qwel",
					UserId: 49,
					Expires: "2022-05-03 12:10:20",
				},
			},
			expectResult: true,
			mockBehaviour: func(args args, expectResult interface{}){
				mock.ExpectExec("INSERT INTO sessions").
					WithArgs(args.session.SessionId, args.session.UserId, args.session.Expires).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
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