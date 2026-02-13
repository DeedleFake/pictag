import { BrowserRouter, useSearchParams } from 'react-router'
import AsyncSelect from 'react-select/async'

import './index.css'

async function loadTags(input: string) {
  return []
}

function Search() {
  const [searchParams, setSearchParams] = useSearchParams()

  return (
    <div className="flex flex-col justify-start items-center">
      <AsyncSelect
        className="w-52"
        isMulti
        defaultOptions
        loadOptions={loadTags}
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
