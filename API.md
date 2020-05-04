# golang-cms 系统

## 全局相关说明 

系统 分为 前台 和 后台。

前台 逻辑 主要服务于 用户，用户交互  和 站点数据展示 数据收集 

后台 逻辑 主要服务于 站点内容 管理  权限管理 数据查询 业务排查

## 用户模块

用户注册，用户登陆 信息管理

### 接口 列表 

|         path         | method |     introduction     | middiewares |
| :------------------: | :----: | :------------------: | :---------: |
|    /user/register    |  post  |     用户注册接口     |             |
|     /user/login      |  post  |     用户登陆接口     |             |
|    /user/profile     |  get   |   用户个人信息接口   |    Auth     |
| /user/:name/profile  |  get   |   用户公共信息接口   |    Limit    |
|      /user/home      |  get   | 用户个人主页信息接口 |    Auth     |
| /user/reset/password |  put   |     用户密码修改     |    Auth     |
|  /user/permissions   |  get   | 用户权限信息获取接口 |    Auth     |
|                      |        |                      |             |
|                      |        |                      |             |



#### 用户模块 接口参数参数说明





## 菜单模块

管理后台菜单 显示 ，前台菜单控制 

### 接口 列表 

|        path        | method |     introduction     | middiewares |
| :----------------: | :----: | :------------------: | :---------: |
|    /menus/lists    |  get   | 获取管理后台菜单接口 |    Auth     |
|     /menus/add     |  post  |     菜单添加接口     |    Auth     |
|   /menus/update    |  put   |     菜单更新接口     |    Auth     |
|   /menus/delete    | delete |     菜单删除接口     |    Auth     |
|    /menus/:name    |  get   |     菜单详情接口     |    Auth     |
|                    |        |                      |             |
|    /menus/:uid     |  get   |   获取对应用户菜单   |    Auth     |
|  /menus/:uid/add   |  post  |     用户菜单添加     |    Auth     |
| /menus/:uid/update |  put   |     用户菜单更新     |    Auth     |



## 权限模块 

后台用户权限管理 ，用户添加，角色管理 

### 接口列表

##### 权限

|             path              | method |   introduction   | middiewares |
| :---------------------------: | :----: | :--------------: | :---------: |
|       /permission/lists       |  get   |   权限列表接口   |    Auth     |
|        /permission/add        |  post  |   权限添加接口   |    Auth     |
|      /permission/update       |  put   |   权限更新接口   |    Auth     |
|      /permission/delete       | delete |   权限删除接口   |    Auth     |
|       /permission/:name       |  get   |   权限详情接口   |    Auth     |
|                               |        |                  |             |
|    /permission/:name/users    |  get   | 对应权限用户列表 |    Auth     |
|  /permission/:name/user/add   |  post  | 对应权限用户添加 |    Auth     |
| /permission/:name/user/delete | delete | 对应权限用户删除 |    Auth     |
|                               |        |                  |             |
|    /permission/:name/roles    |  get   | 对应权限角色列表 |    Auth     |
|  /permission/:name/role/add   |  post  | 对应权限角色添加 |    Auth     |
| /permission/:name/role/delete | delete | 对应角色权限删除 |    Auth     |
|                               |        |                  |             |

##### 角色 

|          path           | method |       introduction       | middiewares |
| :---------------------: | :----: | :----------------------: | :---------: |
|       /role/lists       |  get   |       角色列表接口       |    Auth     |
|        /role/add        |  post  |       角色添加接口       |    Auth     |
|      /role/update       |  put   |       角色更新接口       |    Auth     |
|      /role/delete       | delete |       角色删除接口       |    Auth     |
|       /role/:name       |  get   |       角色详情接口       |    Auth     |
|                         |        |                          |             |
|    /role/:name/users    |  get   | 获取对应角色用户列表接口 |    Auth     |
|  /role/:name/user/add   |  post  |   角色用户用户添加接口   |    Auth     |
| /role/:name/user/delete | delete |   角色用户用户删除接口   |    Auth     |
|                         |        |                          |             |
|                         |        |                          |             |
|                         |        |                          |             |



## Ui 模块

对应接口 逻辑 UI 页面 聚合接口 ， web 路由 

#### 资源路由

|       name        |            router            |                           elements                           |  module  |
| :---------------: | :--------------------------: | :----------------------------------------------------------: | :------: |
|    用户登陆页     |     /web/user/auth#login     |  登陆表单,登陆接口  用户名,手机号,邮箱 密码\|登陆码  验证码  |   user   |
|    用户注册页     |   /web/user/auth#register    | 注册表单,注册接口 用户名,性别,手机号,邮箱,密码,短信验证码,图形验证码 |   user   |
|     用户主页      |       /web/user/:name        |        用户个人信息接口 字段类型样式(头像,姓名,性别)         |   user   |
|  用户修改密码页   |   /web/user/reset/password   | 修改密码表单,修改密码接口 登陆授权信息,xsrf_token,密码,确认密码 |   user   |
|     站点首页      |          /web/index          | 简约布局(首,体,尾)  首页信息聚合接口, 首页文章(内容)列表接口 广告接口 |  index   |
| 内容频道详情(>=1) |      /web/channel/:name      |   子页布局 频道名  频道信息列表接口 广告接口 其他相关接口    | channels |
|                   |                              |                                                              |          |
|  后台用户登陆页   |    /admin/user/auth#login    |  登陆表单,后台登陆接口  用户名,手机号,邮箱 密码 图形验证码   |  admin   |
|     后台首页      |         /admin/home          | 登陆首页布局,菜单接口,实时消息接口,用户权限信息接口,首页信息聚合接口 消息提醒样式 |  admin   |
|    菜单列表页     |         /admin/menus         | 表格布局,分页加载,操作按钮(修改,删除,添加子菜单,详情)  菜单列表接口 |  admin   |
|    菜单详情页     |      /admin/menu/:name       |  样式布局,操作按钮 (修改,删除,添加子菜单) 对应菜单详情接口   |  admin   |
|    菜单添加页     |       /admin/menu#add        | 添加字段表单(菜单名,排序,父级,信息[icon,alias,buttons]) 对应菜单添加接口 |  admin   |
|    菜单修改页     |   /admin/menu/:name#update   |          对应菜单信息，可修改字段表单 菜单更新接口           |  admin   |
|    菜单删除页     |   /admin/menu/:name#delete   |          对应菜单名\|id 菜单删除提醒  菜单删除接口           |  admin   |
|                   |                              |                                                              |          |
|    活动列表页     |        /web/activites        |              展示活动列表 分页加载 活动列表接口              | activity |
|    活动详情页     |     /web/activity/:name      |         活动名,id 活动信息详情接口 活动表单获取接口          | activity |
|  活动表单信息页   |  /web/activity/:name#:form   |            动态表单 动态字段 动态表单字段提交接口            | activity |
|   活动表单导出    | /admin/activity/:name#export |               活动数据导出页 活动数据导出接口                | activity |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |
|                   |                              |                                                              |          |



#### 接口列表



## 内容管理模块 

富文本文章内容编辑 ， 内容列表, 内容类型管理 ，内容发布 内容修改 内容删除 内容下架 内容分类 内容收集

#### 接口列表 

|               path               | method |         introduction         | middiewares |
| :------------------------------: | :----: | :--------------------------: | :---------: |
|          /admin/editor           |  get   |       加载富文本编辑器       |    Auth     |
|                                  |        |                              |             |
|           /admin/posts           |  get   |         内容分页接口         |    Auth     |
|         /admin/post/:id          |  get   |       文章内容详情接口       |    Auth     |
|         /admin/post/add          |  post  |         添加文章接口         |    Auth     |
|        /admin/post/update        |  put   |         文章更新接口         |    Auth     |
|        /admin/post/delete        | delete |         文章删除接口         |    Auth     |
|        /admin/post/types         |  get   |     获取文章类列表型接口     |    Auth     |
|       /admin/post/type/add       |  post  |         添加文章类型         |    Auth     |
|     /admin/post/type/update      |  put   |         文章类型更新         |    Auth     |
|     /admin/post/type/delete      | delete |         文章类型删除         |    Auth     |
|     /admin/post/:type/lists      |  get   | 获取对应类型文章分页列表接口 |    Auth     |
|      /admin/post/:type/add       |  post  |    添加文章分类类型(批量)    |    Auth     |
|     /admin/post/:type/delete     | delete |   删除文章的分类类型(批量)   |    Auth     |
|        /admin/post/:id/on        |  put   |         文章审核发布         |    Auth     |
|       /admin/post/:id/off        |  put   |           文章下架           |    Auth     |
|        /admin/post/search        |  post  |      文章搜索接口(分页)      |    Auth     |
|                                  |        |                              |             |
|     /admin/task/collections      |  get   |     文章采集任务列表接口     |    Auth     |
|    /admin/task/collection/add    |  get   |     文章采集任务添加接口     |    Auth     |
|  /admin/task/collection/update   |  put   |     文章采集任务更新接口     |    Auth     |
|   /admin/task/collection/stop    |  post  |     文章采集任务停止接口     |    Auth     |
|   /admin/task/collection/start   |  post  |     文章采集任务启动接口     |    Auth     |
|   /admin/task/collection/:name   |  get   |   文章采集任务信息详情接口   |    Auth     |
| /admin/collection/:name/contents |  get   | 文章采集任务结果内容列表接口 |    Auth     |
|  /admin/collection/:name/update  |  put   |  文章采集任务结果内容修改口  |    Auth     |
|  /admin/collection/:name/delete  | delete |  文章采集任务结果内容删除口  |    Auth     |
|                                  |        |                              |             |
|                                  |        |                              |             |





##  附件管理模块 

文章图片上传 ， 素材库管理 ，可下载附件服务管理 

#### 接口列表

|          path           | method |           introduction           |  middiewares   |
| :---------------------: | :----: | :------------------------------: | :------------: |
|      /attachments       |  get   |           附件列表接口           |      Auth      |
|   /attachment/upload    |  post  |           附件上传接口           | Auth,Size,Role |
| /attachment/:type/list  |  get   |         对应类型附件列表         |   Auth,Role    |
|     /attachment/:id     |  get   |           附件详情接口           |      Auth      |
|   /attachment/delete    | delete |           附件删除接口           |      Auth      |
|   /attachment/update    |  put   |       附件相关信息更新接口       |      Auth      |
| /attchment/:id/download |  get   |             附件下载             |      Auth      |
|  /attchment/:id/check   |  get   | 附件检查(涉黄，异常，病毒，垃圾) |      Auth      |
| /attachment/:uid/lists  |  get   |      获取某个用户的附件列表      |   Authm,Role   |



##  缓存模块 

管理缓存 ，缓存可视化  优化页面接口加载速度 减少io

#### 接口列表

|        path        | method |       introduction       | middiewares |
| :----------------: | :----: | :----------------------: | :---------: |
|      /caches       |  get   |     缓存策略列表接口     |  Auth,Role  |
|     /cache/:id     |  get   | 获取对应缓存策略信息接口 |  Auth,Role  |
|    /cache/keys     |  get   |    获取缓存键列表接口    |  Auth,Role  |
|     /cache/add     |  post  |     添加缓存策略接口     |  Auth,Role  |
|   /cache/update    |  put   |     修改缓存策略接口     |  Auth,Role  |
|   /cache/delete    | delete |     删除缓存策略接口     |  Auth,Role  |
|  /cache/:key/get   |  get   |    获取对应缓存值接口    |  Auth,Role  |
| /cache/:key/delete | delete |      删除对应缓存值      |  Auth,Role  |
|  /cache/:id/load   |  get   | 运行对应缓存策略加载计划 |  Auth,Role  |

