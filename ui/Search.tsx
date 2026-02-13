import { useSearchParams } from 'react-router'
import AsyncSelect from 'react-select/async'

import { listTags } from './api'

async function loadTags(input: string) {
  const tags = await listTags(input)
  return tags.map((tag) => ({
    label: `${tag.name} (${tag.count})`,
    value: tag.name,
  }))
}

export function Search() {
  const [searchParams, setSearchParams] = useSearchParams()
  const defaultValue = searchParams.getAll('q').map((value) => ({
    label: value,
    value,
  }))

  return (
    <AsyncSelect
      className="w-100"
      isMulti
      defaultOptions
      loadOptions={loadTags}
      defaultValue={defaultValue}
      onChange={(vals) => {
        setSearchParams({ q: vals.map((val) => val.value) })
      }}
    />
  )
}

export default Search
