# Blue Bell Frontend

Blue Bell 是一个社区和帖子管理系统的前端项目，使用 Vue 3 + TypeScript + Tailwind CSS 开发。

## 功能特性

- 用户认证（登录、注册、退出）
- 社区管理（查看社区列表、社区详情）
- 帖子管理（发布、查看、投票）
- 评论管理（发布、查看、投票）
- 响应式设计

## 技术栈

- Vue 3
- TypeScript
- Vue Router
- Pinia
- Axios
- Tailwind CSS
- Vite

## 开发环境要求

- Node.js >= 16.0.0
- npm >= 7.0.0

## 安装和运行

1. 克隆项目

```bash
git clone https://github.com/your-username/bluebell-frontend.git
cd bluebell-frontend
```

2. 安装依赖

```bash
npm install
```

3. 配置环境变量

复制 `.env.example` 文件为 `.env`，并根据需要修改环境变量：

```bash
cp .env.example .env
```

4. 启动开发服务器

```bash
npm run dev
```

5. 构建生产版本

```bash
npm run build
```

## 项目结构

```
bluebell-frontend/
├── src/
│   ├── api/          # API 请求
│   ├── assets/       # 静态资源
│   ├── components/   # 组件
│   ├── router/       # 路由配置
│   ├── stores/       # 状态管理
│   ├── types/        # TypeScript 类型定义
│   ├── views/        # 页面组件
│   ├── App.vue       # 根组件
│   └── main.ts       # 入口文件
├── public/           # 公共资源
├── index.html        # HTML 模板
└── package.json      # 项目配置
```

## API 接口

### 认证

- POST /auth/login - 用户登录
- POST /auth/register - 用户注册
- POST /auth/logout - 用户退出
- GET /auth/me - 获取当前用户信息

### 社区

- GET /communities - 获取社区列表
- GET /communities/:id - 获取社区详情
- POST /communities - 创建社区
- PUT /communities/:id - 更新社区
- DELETE /communities/:id - 删除社区

### 帖子

- GET /posts - 获取帖子列表
- GET /posts/:id - 获取帖子详情
- POST /posts - 创建帖子
- PUT /posts/:id - 更新帖子
- DELETE /posts/:id - 删除帖子
- POST /posts/:id/vote - 投票

### 评论

- GET /posts/:id/comments - 获取评论列表
- POST /comments - 创建评论
- PUT /comments/:id - 更新评论
- DELETE /comments/:id - 删除评论
- POST /comments/:id/vote - 投票

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

[MIT](https://choosealicense.com/licenses/mit/)
