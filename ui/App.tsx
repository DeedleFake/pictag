import { useMemo, useState } from 'react'
import { BrowserRouter, useSearchParams } from 'react-router'

import './index.css'

function splitTags(q: string) {
  return Array.from(new Set(q.split(/\s+/).filter((tag) => tag !== '')))
}

function Search() {
  const [searchParams, setSearchParams] = useSearchParams()
  const tags = useMemo(
    () => Array.from(new Set(searchParams.getAll('q'))),
    [searchParams],
  )
  const [search, setSearch] = useState(tags.join(' '))

  return (
    <div>
      <input
        type="text"
        value={search}
        onChange={(ev) => {
          setSearch(ev.target.value)
          setSearchParams({ q: splitTags(ev.target.value) })
        }}
      />
      <ul>
        {tags.map((tag) => (
          <li key={tag}>{tag}</li>
        ))}
      </ul>
    </div>
  )
}

export function App() {
  return (
    <BrowserRouter>
      <Search />
    </BrowserRouter>
  )
}

export default App
