# agentmem

**หน่วยความจำแบบไม่ใช้เวกเตอร์สำหรับ AI agent**

Agent ลืมข้อมูลข้ามเซสชั่น และที่แย่กว่านั้น เมื่อความทรงจำถูกผูกกับผู้ให้บริการโมเดลหรือ harness เดียว ข้อตกลง การตัดสินใจ และความรู้เกี่ยวกับโปรเจกต์ของคุณจะหายไปพร้อมผลิตภัณฑ์เหล่านั้น

`agentmem` เป็นชั้นความจำแบบปลั๊กอินที่คุณเป็นเจ้าของ ไฟล์ markdown ธรรมดาบนเครื่องของคุณ บริหารจัดการด้วยไบนารี `memory` ตัวเดียว เชื่อมกับ agent ที่รัน shell หรือต่อ MCP ได้ ความทรงจำจึงยังอยู่กับคุณเมื่อเปลี่ยนเครื่องมือ

**ใช้ฟรีทั้งส่วนตัวและเชิงพาณิชย์ ซอร์สเป็นแบบปิด** repository นี้เก็บเว็บไซต์เอกสาร ระบบรายงานปัญหา และไฟล์ไบนารีของเวอร์ชั่นต่าง ๆ

> 🌐 Read in [English](README.md)

## การติดตั้ง

```bash
curl -fsSL https://agentmem.thamtech.co/install | sh
```

หรือดาวน์โหลดไบนารีล่าสุดจาก [GitHub Releases](https://github.com/vtno/agentmem/releases)

### เชื่อมกับ harness

Harness ที่รองรับอย่างเป็นทางการ (skill + enforcement hooks/plugins):

```bash
memory install claude-code
memory install opencode
```

MCP (hooks/plugins ยังถูกติดตั้งบนเป้าหมายทางการ):

```bash
memory install --mcp claude-code
memory install --mcp opencode
```

Harness อื่น ๆ ที่โหลด skill ได้ (portable skill เท่านั้น — ไม่มี hooks หรือ plugins):

```bash
memory install skill
# ระบุรากของ skills เองได้ (ตัวเลือก):
memory install skill --dir /path/to/skills
```

พาธ skill เริ่มต้น: `~/.agents/skills/memory/SKILL.md` (ใช้ร่วมกับ skill root ของ OpenCode)

ตรวจสถานะ: `memory targets`

## ทำไมถึงต้องใช้

- **เป็นของคุณ** — ไฟล์ markdown บนดิสก์ ไม่ถูกล็อกในบัญชีผู้ให้บริการ
- **เชื่อมต่อง่าย** — ติดตั้ง skill หรือ MCP ใช้ร่วมกับ harness ที่คุณใช้อยู่
- **พกพาง่าย** — เปลี่ยน agent แต่ใช้โฟลเดอร์ความทรงจำเดิม
- **ไม่ใช้เวกเตอร์** — ไม่มี embedding ไม่มี daemon ไม่ต้องต่ออินเทอร์เน็ต

## เว็บไซต์เอกสาร

Go server ตัวเล็ก ๆ ฝัง (`//go:embed`) ทั้งไดเรคทอรี [`site/`](site/) ไว้ในไบนารีพกพาตัวเดียว — templates, CSS, JS, assets, language packs และ Markdown

### i18n

ข้อความ UI ทั้งหมดอยู่ในไฟล์ YAML ใต้ `site/lang/` โดย root key คือ **ชื่อหน้า** (`common`, `index`, `install`, `commands`, `integrations`, `how_it_works`)

```yaml
# site/lang/en.yaml
common:
  nav:
    install: Install
index:
  title: "agentmem — vectorless memory for AI agents"
  tagline: Vectorless memory for AI agents
```

Templates เรียกใช้ `{{.T "key"}}` (หาในหน้าปัจจุบันก่อน แล้วจึงไป `common`) หรือใช้พาธเต็ม `{{.T "common.nav.install"}}`

การเลือกภาษา (ใช้ค่าที่ตรงก่อน):

1. `?lang=en`
2. Cookie `lang=en`
3. `Accept-Language`
4. ค่าเริ่มต้น `en`

เพิ่มภาษาได้โดยคัดลอก `site/lang/en.yaml` → `site/lang/<code>.yaml` แล้วแปลค่าต่าง ๆ

| เส้นทาง | |
|-------|--|
| `/` | หน้าแรก |
| `/docs/install` | ไบนารี + การเชื่อม harness |
| `/docs/commands` | อ้างอิง CLI |
| `/docs/integrations` | Claude Code, OpenCode, harness อื่น ๆ |
| `/docs/how-it-works` | การจัดเก็บ ดัชนี UI ในเครื่อง |
| `/*.md` / `/docs/*.md` | หน้าเดียวกันในรูปแบบ Markdown (`Content-Type: text/markdown`) |

### Makefile

```bash
make build                 # → ./agentmem-site
make install               # → ~/.local/bin/agentmem-site
make run                   # go run . -addr :5555  (ADDR=:8080 เพื่อ override)
make test
make fmt
VERSION=0.1.0 make build   # ฝังเวอร์ชั่นลงในไบนารี
./agentmem-site -version
```

### สร้าง / รันเอง

```bash
# พรีวิว (สร้าง embed ใหม่จาก site/)
go run . -addr :5555

# ไบนารีพกพา (ไม่ต้องมีไดเรคทอรี site/ ตอนรัน)
make build
./agentmem-site -addr :5555
# ส่งเฉพาะ agentmem-site ไปที่เซิร์ฟเวอร์
```

โครงสร้าง: `site/templates/` (`layout`, `header`, `footer` + เนื้อหาแต่ละหน้า)  
Static: `site/css`, `site/js`, `site/assets` · Markdown: `site/index.md`, `site/docs/*.md`  
การเปลี่ยนเนื้อหาต้อง build ใหม่เพื่อฝังกลับเข้าไป

## คำสั่งหลัก

```bash
memory list [prefix]          # แสดงสารบัญความจำ ใส่ prefix เพื่อกรองได้
memory show <name>            # อ่านความจำหนึ่งไฟล์
memory search <query>         # ค้นหา (grep) ในความทรงจำ
memory add <text>             # เขียนโน้ตแบบไว ๆ เพิ่มเข้าไป
memory add --file <name> --summary "[tag] desc" <text>
memory edit <name>            # เปิดใน $EDITOR
memory rm <name>              # ลบความจำ
memory reindex                # สร้าง MEMORY.md ใหม่
memory serve                  # MCP server บน stdio
memory ui                     # web UI บนเครื่อง (loopback)
memory install <harness>      # skill + hooks/plugins (claude-code, opencode)
memory install --mcp <harness>
memory install skill [--dir <dir>]
memory uninstall …
memory targets                # สถานะการติดตั้ง harness + skill
```

## หลักการทำงาน

ทุกอย่างอยู่ใน `~/.agentmem/` โดยค่าเริ่มต้น (เปลี่ยนด้วย `AGENTMEM_DIR`) แต่ละความจำเป็นไฟล์ markdown `MEMORY.md` คือดัชนีที่สร้างขึ้นจากบรรทัดสรุปของแต่ละไฟล์ ไฟล์คือแหล่งความจริงเดียว

## Harness ที่รองรับ

| เป้าหมาย | สิ่งที่ได้ |
|--------|----------------|
| **Claude Code** (ทางการ) | Skill + SessionStart hook · MCP เป็นตัวเลือก |
| **OpenCode** (ทางการ) | Skill + plugin · MCP เป็นตัวเลือก |
| **skill** (พกพา) | Skill อย่างเดียวใต้ `~/.agents/skills` — ไม่มี enforcement |

## รายงานปัญหาและเวอร์ชั่น

- ปัญหาและฟีเจอร์: [GitHub Issues](https://github.com/vtno/agentmem/issues)
- ไบนารี: [GitHub Releases](https://github.com/vtno/agentmem/releases)

## สัญญาอนุญาต

ดู [LICENSE](LICENSE) ใช้ฟรี แต่ซอร์สไม่เปิด
