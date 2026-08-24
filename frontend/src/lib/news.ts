export interface NewsItem {
  id: string
  clubId: string
  clubName: string
  clubTla: string
  title: string
  sourceName: string
  publishedAt: string
  linkUrl: string
}

export interface NewsFeed {
  id: string
  clubId: string
  clubName: string
  url: string
  sourceName: string
  enabled: boolean
  updatedAt: string
}
