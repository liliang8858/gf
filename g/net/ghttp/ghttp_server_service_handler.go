// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 服务注册.

package ghttp

import (
    "errors"
    "strings"
    "gitee.com/johng/gf/g/util/gstr"
)

// 注意该方法是直接绑定函数的内存地址，执行的时候直接执行该方法，不会存在初始化新的控制器逻辑
func (s *Server) BindHandler(pattern string, handler HandlerFunc) error {
    return s.bindHandlerItem(pattern, &handlerItem{
        ctype : nil,
        fname : "",
        faddr : handler,
    })
}

// 绑定URI到操作函数/方法
// pattern的格式形如：/user/list, put:/user, delete:/user, post:/user@johng.cn
// 支持RESTful的请求格式，具体业务逻辑由绑定的处理方法来执行
func (s *Server) bindHandlerItem(pattern string, item *handlerItem) error {
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server handlers cannot be changed while running")
    }
    return s.setHandler(pattern, item)
}

// 通过映射数组绑定URI到操作函数/方法
func (s *Server) bindHandlerByMap(m handlerMap) error {
    for p, h := range m {
        if err := s.bindHandlerItem(p, h); err != nil {
            return err
        }
    }
    return nil
}

// 将内置的名称按照设定的规则合并到pattern中，内置名称按照{.xxx}规则命名。
// 规则1：pattern中的URI包含{.struct}关键字，则替换该关键字为结构体名称；
// 规则1：pattern中的URI包含{.method}关键字，则替换该关键字为方法名称；
// 规则2：如果不满足规则1，那么直接将防发明附加到pattern中的URI后面；
func (s *Server) mergeBuildInNameToPattern(pattern string, structName, methodName string) string {
    structName = s.nameToUrlPart(structName)
    methodName = s.nameToUrlPart(methodName)
    pattern    = strings.Replace(pattern, "{.struct}", structName, -1)
    if strings.Index(pattern, "{.method}") != -1 {
        return strings.Replace(pattern, "{.method}", methodName, -1)
    }
    // 检测域名后缀
    array := strings.Split(pattern, "@")
    // 分离URI(其实可能包含HTTP Method)
    uri := array[0]
    uri  = strings.TrimRight(uri, "/") + "/" + methodName
    // 加上指定域名后缀
    if len(array) > 1 {
        return uri + "@" + array[1]
    }
    return uri
}

// 将给定的名称转换为URL规范格式。
// 规范1: 全部转换为小写；
// 规范2: 方法名中间存在大写字母，转换为小写URI地址以“-”号链接每个单词；
func (s *Server) nameToUrlPart(name string) string {
    part := ""
    for i := 0; i < len(name); i++ {
        if i > 0 && gstr.IsLetterUpper(name[i]) {
            part += "-"
        }
        part += string(name[i])
    }
    return strings.ToLower(part)
}