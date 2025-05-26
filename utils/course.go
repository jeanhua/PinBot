// 查课
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Course struct{}

type SearchResult struct {
	Code int `json:"code"`
	Data []struct {
		Kid          int    `json:"kid"`          //id
		Kch          int    `json:"kch"`          //课程号
		Kxh          string `json:"kxh"`          //课序号
		TeachersName string `json:"tname"`        //教师名
		CourseName   string `json:"kname"`        //课程名
		ExamTypeName string `json:"examTypeName"` //考察类型
	} `json:"data"`
}

type CourseDetail struct {
	Code int `json:"code"`
	Data struct {
		Kid          int     `json:"kid"`          //id
		Kch          int     `json:"kch"`          //课程号
		Kxh          string  `json:"kxh"`          //课序号
		TeachersName string  `json:"tname"`        //教师名
		CourseName   string  `json:"kname"`        //课程名
		ExamTypeName string  `json:"examTypeName"` //考察类型
		Credit       float32 `json:"credit"`       //学分

		Count        int     `json:"count"` //统计人数
		Average      float32 `json:"avg"`   //平均分
		Max          float32 `json:"max"`   //最高分
		Min          float32 `json:"min"`   //最低分
		A_levelCount int     `json:"a"`     //90~100分段人数
		B_levelCount int     `json:"b"`     //80~89分段人数
		C_levelCount int     `json:"c"`     //70~79分段人数
		D_levelCount int     `json:"d"`     //60~69分段人数
		E_levelCount int     `json:"e"`     //0~59分段人数

		History []struct {
			ExamTime     int     `json:"examTime"` //考试时间
			Count        int     `json:"count"`    //统计人数
			Average      float32 `json:"avg"`      //平均分
			Max          float32 `json:"max"`      //最高分
			Min          float32 `json:"min"`      //最低分
			A_levelCount int     `json:"a"`        //90~100分段人数
			B_levelCount int     `json:"b"`        //80~89分段人数
			C_levelCount int     `json:"c"`        //70~79分段人数
			D_levelCount int     `json:"d"`        //60~69分段人数
			E_levelCount int     `json:"e"`        //0~59分段人数
		} `json:"history"` //历史分数
	} `json:"data"`
}

func (Course) Search(keyWords string, page int) (*SearchResult, error) {
	client := http.Client{}
	link := fmt.Sprintf("https://duomi.chenyipeng.com/pennisetum/scu/score/search?kname=%s&page=%d", url.QueryEscape(keyWords), page)
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result SearchResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (Course) GetDetail(kid int) (*CourseDetail, error) {
	client := http.Client{}
	link := fmt.Sprintf("https://duomi.chenyipeng.com/pennisetum/scu/score/getDetail?kid=%d", kid)
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result CourseDetail
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
