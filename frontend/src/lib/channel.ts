const iconify = (name: string) => `https://api.iconify.design/simple-icons/${name}.svg`
const commons = (name: string) => `https://commons.wikimedia.org/wiki/Special:Redirect/file/${name}`

const channelIcons: Array<[RegExp, string]> = [
  [/^canal\+/, commons('CanalPlus.svg')],
  [/^dazn\b/, iconify('dazn')],
  [/^be ?in sports\b/, commons('BeIN-Sports-Logo.svg')],
  [/^disney\+(?:\s|$)/, commons('Disney%2B_logo.svg')],
  [/^(amazon )?prime video\b/, iconify('primevideo')],
  [/^apple tv(?:\+)?(?:\s|$)/, iconify('appletv')],
  [/^youtube\b/, iconify('youtube')],
]

export function channelIconUrl(channel: string): string | null {
  const normalized = channel.trim().toLocaleLowerCase('fr')
  const match = channelIcons.find(([pattern]) => pattern.test(normalized))
  return match?.[1] ?? null
}
