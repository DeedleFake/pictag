export interface Tag {
  name: string
  count: number
}

export async function listTags(name: string): Promise<Tag[]> {
  type Result = { Name: string; Count: number }

  const params = new URLSearchParams()
  params.set('q', name)
  const rsp = await fetch(`/api/list_tags?${params}`)
  const tags = await rsp.json()

  return tags.map((tag: Result) => ({
    name: tag.Name,
    count: tag.Count,
  }))
}
