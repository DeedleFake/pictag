import { BrowserRouter, useSearchParams } from 'react-router'
import AsyncSelect from 'react-select/async'

import { listTags } from './api'
import './index.css'

async function loadTags(input: string) {
  const tags = await listTags(input)
  return tags.map((tag) => ({
    label: `${tag.name} (${tag.count})`,
    value: tag.name,
  }))
}

function Search() {
  const [searchParams, setSearchParams] = useSearchParams()

  return (
    <div className="flex flex-col justify-start items-center">
      <AsyncSelect
        className="w-100"
        isMulti
        defaultOptions
        loadOptions={loadTags}
        onChange={(vals) => {
          setSearchParams({ q: vals.map((val) => val.value) })
        }}
      />
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
