#!/usr/bin/env node
import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js'
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js'
import { AiAgentLogsClient } from './client.js'
import { registerTools } from './tools.js'

async function main(): Promise<void> {
  const client = new AiAgentLogsClient()
  const server = new McpServer({
    name: 'ai-agent-logs-mcp',
    version: '1.0.0',
  })

  registerTools(server, client)

  const transport = new StdioServerTransport()
  await server.connect(transport)
  console.error(`ai-agent-logs-mcp ready (API ${client.baseUrl})`)
}

main().catch((err) => {
  console.error(err)
  process.exit(1)
})
