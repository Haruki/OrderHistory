import { useState } from 'react';
import reactLogo from './assets/react.svg';
import './App.css';

/*Root Component */
const App = () => {
  const title = 'wurstbrot';
  const arr = [
    { id: '1', value: 'butter' },
    { id: '2', value: 'keks' },
    { id: '3', value: 'F1' },
  ];
  return (
    <div>
      <h1>Vite + {title}</h1>
      <List list={arr} />
      <Search />
    </div>
  );
};

/* List Component */
const List = (props) => {
  return (
    <ol>
      {props.list.map((x) => (
        <Item key={x.id} item={x} />
      ))}
    </ol>
  );
};

/*Item Component*/
const Item = (props) => <li key={props.item.id}>{props.item.value}</li>;

/* Search Component */
const Search = () => {
  const handleChange = (event) => {
    console.log(event);
    console.log(event.target.value);
  };
  return (
    <div>
      <label htmlFor='search'>Search: </label>
      <input type='text' id='search' onChange={handleChange} />
    </div>
  );
};

export default App;
