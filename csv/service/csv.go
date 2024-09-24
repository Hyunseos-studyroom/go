package service

import (
	"encoding/csv"
	"github.com/csv/types"
	"mime/multipart"
	"strconv"
)

type CSVService struct{}

func (s *CSVService) ReadCSV(file *multipart.FileHeader) ([]types.CSV, error) {
	// 파일을 열기
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// CSV 리더 생성
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// CSV 데이터가 비어있는지 확인
	if len(records) == 0 {
		return nil, nil
	}

	var csvData []types.CSV
	for _, record := range records[1:] { // 헤더를 제외하고 읽기
		if len(record) < 3 { // 각 레코드가 적어도 3개의 요소를 가져야 함
			continue
		}
		age, err := strconv.Atoi(record[1]) // 오류 처리 추가
		if err != nil {
			continue // 오류가 발생하면 해당 레코드는 무시
		}
		csvData = append(csvData, types.CSV{
			Name: record[0],
			Age:  age,
			Job:  record[2],
		})
	}

	return csvData, nil
}
