import { useState } from 'react'
import reactLogo from './assets/react.svg'
import './App.css'

function App() {
   
  const title = "wurstbrot"

  return (
      <div>
      <h1>Vite + {title}</h1>

      <label htmlFor="search">Search: </label>
      <input type="text" id="search" />
    </div>
  )
}

export default App
