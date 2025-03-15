# Danta-Auto-Tool

# 简介

旦挞自动化工具，目前支持茶楼宣传位的申请、审批、上线、通知自动化。更多功能有待增加。该工具利用了飞书的开放平台。

## 使用方式

推荐使用 Dockerfile 来运行该项目，并且需要设置好环境变量。以下是需要设置的环境变量及其含义：

| 环境变量                          | 含义                                      |
|-----------------------------------|-------------------------------------------|
| LARK_USER_ACCESS_TOKEN            | 飞书用户访问令牌                          |
| LARK_APP_ID                       | 飞书应用的 APP ID                         |
| LARK_APP_SECRET                   | 飞书应用的 APP Secret                     |
| LARK_BANNER_APPROVE_CARD_ID       | 飞书应用中的审批卡片 ID                   |
| LARK_BANNER_BITABLE_APP_TOKEN     | Banner 宣传位的多维表格的 APP Token       |
| LARK_BANNER_BITABLE_APPLICATION_TABLE_ID | Banner 宣传位的申请表 Table ID       |
| LARK_BANNER_BITABLE_USAGE_TABLE_ID | Banner 宣传位的使用记录表 Table ID       |
| LARK_BANNER_APPROVE_GROUP_ID      | Banner 宣传位的审批群 ID                  |
| DANTA_DEV_EMAIL                   | Danta 开发者邮箱（暂时没有用到）          |
| GITHUB_PERSONAL_ACCESS_TOKEN      | Github 个人访问令牌                       |
| GITHUB_DANXI_REPO_OWNER           | Github 仓库的 owner                       |
| GITHUB_DANXI_REPO_NAME            | Github 仓库的 name                        |
| GITHUB_DANXI_REPO_APP_CONFIG_PATH | Github 仓库的 Banner 配置文件路径         |

使用 Dockerfile 运行该项目的示例：

```shell
docker build -t danta-auto-tool .
```
```shell
# Run the Docker container with environment variables from .env file
docker run --env-file .env danta-auto-tool
```

## 技术方案

更多技术细节请参考：[技术方案](https://danxi-dev.feishu.cn/wiki/A5mjwoQrWixsvKk73itc2Eoinkd)