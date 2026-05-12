import { createI18n } from 'vue-i18n'
import zhCN from './locales/zh-CN'
import enUS from './locales/en-US'

const savedLang = typeof localStorage !== 'undefined'
  ? (localStorage.getItem('opentools-lang') || 'zh-CN')
  : 'zh-CN'

export const i18n = createI18n({
  locale: savedLang,
  fallbackLocale: 'en-US',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS,
  },
})

export function setLanguage(lang: string) {
  ;(i18n.global.locale as any) = lang
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('opentools-lang', lang)
  }
}

export function getStoredTheme() {
  return typeof localStorage !== 'undefined' ? (localStorage.getItem('opentools-theme') || 'dark') : 'dark'
}

export function getStoredColor() {
  return typeof localStorage !== 'undefined' ? (localStorage.getItem('opentools-color') || 'blue') : 'blue'
}

export function setTheme(theme: string) {
  document.documentElement.setAttribute('data-theme', theme)
  if (typeof localStorage !== 'undefined') localStorage.setItem('opentools-theme', theme)
}

export function setColor(color: string) {
  document.documentElement.className = `theme-${color}`
  if (typeof localStorage !== 'undefined') localStorage.setItem('opentools-color', color)
}

export function initTheme() {
  setTheme(getStoredTheme())
  setColor(getStoredColor())
}
