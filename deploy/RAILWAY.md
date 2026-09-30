# Railway deployment

This fork supports Railway's runtime `PORT` automatically when `MCP_HTTP_ADDR`
is not explicitly set.

## 1. Create the Railway service

Deploy this GitHub repository with the included `Dockerfile`.

Generate a public Railway domain for the service before setting
`MCP_AUTH_ISSUER_URL`.

## 2. Attach persistent storage

Attach one Railway Volume to the service and mount it at:

```text
/data
```

The Docker image normally runs as a non-root user, while Railway volumes are
mounted as root. Set this Railway platform variable so the service can write to
the mounted volume:

```text
RAILWAY_RUN_UID=0
```

Telegram/OAuth session data should then be stored in:

```text
MCP_AUTH_SESSION_DIR=/data/sessions
```

The MCP server encrypts Telegram session blobs before they are written to the
session backend.

## 3. Required service variables

Set these in Railway's Variables tab:

```text
MCP_TRANSPORT=http
MCP_TELEGRAM_API_ID=<Telegram API ID>
MCP_TELEGRAM_API_HASH=<Telegram API hash>
MCP_AUTH_ISSUER_URL=https://<your-railway-domain>
MCP_AUTH_ALLOWED_USERS=<your numeric Telegram user id>
MCP_AUTH_TOKEN_KEYS=<base64 encoded random 32-byte key>
MCP_AUTH_SESSION_DIR=/data/sessions
RAILWAY_RUN_UID=0
```

Do not set `MCP_HTTP_ADDR` on Railway. Railway injects `PORT`, and this fork
binds to `:$PORT` automatically.

For a single-user deployment, use the exact numeric Telegram user ID rather
than `*` once it is known.

Generate the token/session master key locally, for example:

```bash
openssl rand -base64 32
```

Store it only as a Railway secret. Do not commit it to Git.

## 4. Summarisation

The normal Telegram tools do not require a second LLM. Do not configure
Anthropic or Gemini API keys unless you explicitly want `SummarizeChat` to send
chat content to those providers.

The default `sampling` provider does not make direct Anthropic/Gemini requests
from this server.

## 5. OAuth redirect

The server uses embedded OAuth 2.1 with Dynamic Client Registration and PKCE.
If the MCP client requires a redirect URI that is not accepted automatically,
add its exact HTTPS callback URI to:

```text
MCP_AUTH_ALLOWED_REDIRECTS=<exact callback URI>
```

## 6. First connection

After the Railway deployment is healthy, add the remote MCP endpoint to the
client. Complete the OAuth flow, scan the Telegram QR code in Telegram's
Devices screen, and approve the new Telegram session.

After login, verify the account with `GetMe` before using write tools.

## Security notes

- Keep `MCP_TELEGRAM_API_HASH` and `MCP_AUTH_TOKEN_KEYS` only in Railway secrets.
- Keep the session volume attached and backed up as appropriate.
- Restrict `MCP_AUTH_ALLOWED_USERS` to your Telegram numeric user ID.
- If the Railway deployment or token key is compromised, revoke the Telegram
  device session in Telegram and rotate `MCP_AUTH_TOKEN_KEYS`.
- Message/search text is redacted from normal MCP request logs; identifiers and
  operational metadata may still appear in logs.
