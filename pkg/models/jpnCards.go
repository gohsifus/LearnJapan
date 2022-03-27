package models

import(
	"fmt"
)

func GetById(id string) JpnCards{
	rows, err := DB.Query("select * from Cards where id = ?", id)
	if err != nil{
		panic(fmt.Sprintf("Ошибка выбора даных: %s", err))
	}
	defer rows.Close()

	card := JpnCards{}
	if(rows.Next()){
		rows.Scan(&card.id, &card.inJapan, &card.inRussian, &card.mark, &card.dateAdd)
	}

	fmt.Println(id)
	fmt.Println(card)

	return card
}

func GetList() []JpnCards{
	rows, err := DB.Query("select * from Cards")
	if err != nil{
		panic(fmt.Sprintf("Ошибка выбора даных: %s", err))
	}
	defer rows.Close()

	cards := []JpnCards{}

	for rows.Next(){
		card := JpnCards{}
		err = rows.Scan(&card.id, &card.inJapan, &card.inRussian, &card.mark, &card.dateAdd)
		cards = append(cards, card)

		fmt.Println(card)
	}

	return cards
}
