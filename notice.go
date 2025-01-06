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

func (s *Student) GetNoticeInfo() (list []*NoticeInfo, err error) {
	// 获取通知公告
	// 1. 获取通知公告页面
	res, err := s.PostWithIdentifier(constants.NoticeInfoQueryURL, map[string]string{})
	if err != nil {
		return nil, err
	}
	// 2. 解析页面
	list, err = parseNoticeInfo(res)
	if err != nil {
		return nil, err
	}
	// 3. 返回结果
	return list, nil
}

func parseNoticeInfo(doc *html.Node) ([]*NoticeInfo, error) {
	// 解析通知公告页面
	var list []*NoticeInfo

	// 修正 XPath 表达式
	sel := htmlquery.FindOne(doc, "//div[@class='box-gl clearfix']")
	if sel == nil {
		return nil, fmt.Errorf("cannot find the notice list")
	}

	// 查找所有的 <li> 元素
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
