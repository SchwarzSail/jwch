/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package jwch

import (
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"

	"github.com/west2-online/jwch/constants"
)

func (s *Student) GetNoticeInfo(req *NoticeInfoReq) (list []*NoticeInfo, err error) {
	// 获取通知公告页面的总页数
	res, err := s.PostWithIdentifier(constants.NoticeInfoQueryURL, map[string]string{})
	if err != nil {
		return nil, err
	}
	// 首页直接爬取
	if req.PageNum == 1 {
		list, err = parseNoticeInfo(res)
		if err != nil {
			return nil, err
		}
		return list, nil
	}
	// 分页需要根据页数计算 url
	lastPageNum, err := getTotalPages(res)
	if err != nil {
		return nil, err
	}
	// 判断是否超出总页数
	if req.PageNum > lastPageNum {
		return nil, fmt.Errorf("超出总页数")
	}
	// 根据总页数计算 url
	num := lastPageNum - req.PageNum + 1
	url := fmt.Sprintf("https://jwch.fzu.edu.cn/jxtz/%d.htm", num)
	doc, err := s.PostWithIdentifier(url, map[string]string{})
	if err != nil {
		return nil, err
	}
	list, err = parseNoticeInfo(doc)
	if err != nil {
		return nil, err
	}
	// 3. 返回结果
	return list, nil
}

// 获取当前页面的所有数据信息
func parseNoticeInfo(doc *html.Node) ([]*NoticeInfo, error) {
	// 解析通知公告页面
	var list []*NoticeInfo

	sel := htmlquery.FindOne(doc, "//div[@class='box-gl clearfix']")
	if sel == nil {
		return nil, fmt.Errorf("cannot find the notice list")
	}

	rows := htmlquery.Find(sel, ".//ul[@class='list-gl']/li")

	for _, row := range rows {
		// 提取日期
		dateNode := htmlquery.FindOne(row, ".//span[@class='doclist_time']")
		date := strings.TrimSpace(htmlquery.InnerText(dateNode))

		// 提取标题
		titleNode := htmlquery.FindOne(row, ".//a")

		title := strings.TrimSpace(htmlquery.SelectAttr(titleNode, "title"))

		// 提取 URL
		url := strings.TrimSpace(htmlquery.SelectAttr(titleNode, "href"))
		url = constants.JwchNoticeURLPrefix + url

		noticeInfo := &NoticeInfo{
			Title: title,
			URL:   url,
			Date:  date,
		}
		list = append(list, noticeInfo)
	}

	return list, nil
}

// 获取总页数
func getTotalPages(doc *html.Node) (int, error) {
	totalPagesNode := htmlquery.FindOne(doc, "//span[@class='p_pages']//a[@href='jxtz/1.htm']")
	if totalPagesNode == nil {
		return 0, fmt.Errorf("未找到总页数")
	}

	totalPagesStr := htmlquery.InnerText(totalPagesNode)
	var totalPages int
	_, err := fmt.Sscanf(totalPagesStr, "%d", &totalPages)
	if err != nil {
		return 0, fmt.Errorf("解析总页数失败: %v", err)
	}
	return totalPages, nil
}
