export interface MatchCmd {
  type: 'regex' | 'over' | 'img' | 'files' | 'window'
  label: string
  match?: string
  minLength?: number
  exclude?: string
  fileType?: 'file' | 'directory'
  extensions?: string[]
}

export interface CommandItem {
  id?: string
  name: string
  path: string
  icon?: string
  type?: string
  subType?: string
  pluginName?: string
  featureCode?: string
  cmdType?: string
  aliases?: string[]
}

interface MatchContext {
  searchQuery: string
  pastedImage: string | null
  pastedFiles: { name: string; path: string; isDirectory: boolean }[] | null
  pastedText: string | null
  activeWindowApp?: string
}

export function matchCommand(
  cmd: MatchCmd,
  ctx: MatchContext
): boolean {
  switch (cmd.type) {
    case 'regex': {
      if (!cmd.match) return false
      const q = ctx.searchQuery
      if (cmd.minLength && q.length < cmd.minLength) return false
      try {
        const re = new RegExp(cmd.match, 'im')
        return re.test(q)
      } catch { return false }
    }
    case 'over': {
      const q = ctx.searchQuery
      if (cmd.minLength && q.length < cmd.minLength) return false
      if (cmd.exclude) {
        try { if (new RegExp(cmd.exclude, 'im').test(q)) return false } catch { }
      }
      return q.length > 0
    }
    case 'img': {
      return ctx.pastedImage !== null
    }
    case 'files': {
      if (!ctx.pastedFiles || ctx.pastedFiles.length === 0) return false
      if (cmd.minLength && ctx.pastedFiles.length < cmd.minLength) return false
      return true
    }
    case 'window': {
      if (!ctx.activeWindowApp) return false
      return true
    }
    default:
      return false
  }
}

export function getMatchLabel(cmd: MatchCmd, ctx: MatchContext): string {
  switch (cmd.type) {
    case 'img': return cmd.label || 'Search image'
    case 'files': return cmd.label || `Search ${ctx.pastedFiles?.length || ''} files`
    case 'window': return cmd.label || `Active: ${ctx.activeWindowApp || ''}`
    default: return cmd.label
  }
}
