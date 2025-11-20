这份文档采用了 **Golang (后端) + Tailwind CSS (样式) + Supabase (数据库 & Auth)** 的架构，这是一种开发速度快、性能高且维护成本低的组合。

-----

# 📝 项目技术开发文档：英语语法通 (English Grammar Master)

## 1\. 项目概况 (Project Overview)

  * **项目名称**：English Grammar Master (EGM)
  * **目标用户**：以中文为母语的英语学习者。
  * **核心功能**：分级语法教学（初/中/高）、配套练习题、用户进度追踪。
  * **技术栈**：
      * **后端**：Golang (使用 Gin框架)。
      * **前端渲染**：Golang `html/template` (最后打包后，HTML和其他静态资源需要打包进二进制文件,方便部署) 。
      * **样式**：Tailwind CSS (通过 CDN)。
      * **数据库 & 认证**：Supabase (PostgreSQL, Auth, Row Level Security)。

-----

## 2\. 数据库设计 (Database Schema - Supabase)

*请 AI 根据以下 Schema 在 Supabase 中创建表结构。*

### 2.1 Enums (枚举)

  * `difficulty_level`: 'beginner', 'intermediate', 'advanced'
  * `question_type`: 'choice' (选择题), 'fill\_blank' (填空题)

### 2.2 Tables (数据表)

1.  **`profiles`** (用户档案，与 `auth.users` 联动)

      * `id`: UUID (Primary Key, references `auth.users.id`)
      * `username`: Text
      * `avatar_url`: Text
      * `current_level`: `difficulty_level` (Default: 'beginner')
      * `created_at`: Timestamp

2.  **`courses`** (课程/语法点)

      * `id`: Serial (Primary Key)
      * `title`: Text (e.g., "一般现在时 - Present Simple")
      * `level`: `difficulty_level`
      * `content_md`: Text (语法讲解，支持 Markdown 格式)
      * `order_index`: Integer (用于排序学习顺序)
      * `slug`: Text (用于 URL, e.g., "present-simple")

3.  **`exercises`** (练习题)

      * `id`: Serial (Primary Key)
      * `course_id`: Integer (FK to `courses.id`)
      * `question`: Text (题目内容)
      * `options`: JSONB (如果是选择题，存储选项数组)
      * `correct_answer`: Text (正确答案)
      * `explanation`: Text (答案解析，中文)
      * `type`: `question_type`

4.  **`user_progress`** (学习进度)

      * `id`: Serial (Primary Key)
      * `user_id`: UUID (FK to `profiles.id`)
      * `course_id`: Integer (FK to `courses.id`)
      * `is_completed`: Boolean
      * `score`: Integer (最近一次得分)
      * `completed_at`: Timestamp

-----

## 3\. 路由与 API 设计 (Routes Structure)

### 3.1 公开路由 (Public)

  * `GET /`: 首页 (Landing Page)，展示介绍和等级入口。
  * `GET /login`: 登录页面 (集成 Supabase Auth UI)。
  * `GET /signup`: 注册页面。

### 3.2 保护路由 (Protected - 需要中间件验证 Auth)

  * **仪表盘**:
      * `GET /dashboard`: 展示用户当前等级、学习进度概览。
  * **课程学习**:
      * `GET /learn/:level`: 列出该等级下的所有语法课程列表 (标记已学/未学)。
      * `GET /learn/:level/:slug`: 具体的语法学习页面 (展示 `content_md`)。
  * **练习与提交**:
      * `GET /learn/:level/:slug/quiz`: 获取该课程的练习题页面。
      * `POST /api/quiz/submit`: 提交答案，后端计算分数，更新 `user_progress`，返回结果和解析。

-----

## 4\. 页面功能详情 (Feature Specifications)

### 4.1 首页 & 导航

  * 导航栏：Logo，等级切换（初/中/高），用户头像/登录按钮。
  * Tailwind 风格：简洁、清新（建议使用 Teal 或 Indigo 色系）。

### 4.2 语法学习页

  * 左侧/顶部：面包屑导航 (首页 \> 初级 \> 一般现在时)。
  * 中间：语法内容区域。需将数据库中的 Markdown 渲染为 HTML。
  * 底部：**"开始练习"** 按钮。

### 4.3 练习页 (核心功能)

  * 展示 5-10 个题目。
  * **UI 交互**：
      * 用户选择或输入答案。
      * 点击“提交”后，通过 AJAX (或 HTMX) 发送请求。
      * 界面即时反馈：正确的显示绿色，错误的显示红色并展示 `explanation` (解析)。
      * 如果是满分或及格，显示“下一课”按钮。

### 4.4 用户系统

  * 利用 Supabase GoTrue (Auth) 处理注册/登录。
  * 后端中间件 (Middleware) 解析 Cookie/Token，确认用户身份。

-----

## 5\. 目录结构建议 (Go Project Layout)

```text
/cmd/server/main.go      # 入口文件
/internal
    /models              # 数据库结构体
    /handlers            # HTTP 处理逻辑 (Dashboard, Quiz, etc.)
    /middleware          # Auth 中间件
    /database            # Supabase 连接初始化
    /services            # 业务逻辑 (如计算分数)
/views                   # HTML 模板 (Go Templates)
    /layouts             # base.html (包含 Tailwind CDN)
    /pages               # index.html, quiz.html
/static                  # 静态资源 (images, css)
```

-----

## 6\. 如何让 AI 帮你开发 (Prompting Strategy)

为了防止 AI 一次性生成太多代码导致混乱，建议你**分阶段**向 AI 发送指令：

### 第一阶段：环境与基础

> "我正在做一个基于 Golang (Gin) + Tailwind + Supabase 的英语学习网站。请帮我初始化项目结构，编写 `main.go`，并设置好 Supabase 的连接代码。请使用 Go Modules。"

### 第二阶段：数据库与模型

> "这是我的数据库 Schema 设计（复制上面的第2部分）。请帮我编写 Supabase 的 SQL 建表语句，并在 Golang 的 `/internal/models` 目录下创建对应的 Struct。"

### 第三阶段：页面与路由

> "请帮我创建基础的 HTML 模板（使用 Tailwind CSS 美化）。我需要一个 `base.html` 布局文件和一个 `index.html` 首页。然后在 Gin 中配置路由来渲染这些页面。"

### 第四阶段：核心业务（课程与练习）

> "现在我们要实现‘课程详情页’和‘练习题’功能。
>
> 1.  请写一个 Handler 从数据库获取 Markdown 内容并在前端渲染。
> 2.  请写一个练习题页面，支持表单提交。
> 3.  请写一个 POST 接口，接收用户答案，比对数据库中的正确答案，并返回 JSON 格式的成绩和解析。"

-----

### 💡 建议的下一步

**你可以直接复制以下这段话发给 AI 开始第一步：**

> “你好，我要开发一个基于 Golang (Gin框架) + Tailwind CSS + Supabase 的英语语法学习网站。
>
> 这是我的**数据库设计**：
> (此处粘贴上面的第2部分 '数据库设计')
>
> 请先帮我生成一份 **Supabase 的 SQL 建表脚本**，确保包含初级、中级、高级的枚举类型定义。生成的 SQL 代码需要可以直接在 Supabase 的 SQL Editor 中运行。”
