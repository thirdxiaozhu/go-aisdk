/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-11 15:58:42
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-11 16:00:09
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package main

import (
	"fmt"
	"github.com/liusuxian/aisdk"
)

func main() {
	factory := aisdk.NewProviderFactory()
	fmt.Println("factory: ", factory)
}
