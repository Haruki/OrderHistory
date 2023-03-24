import { useState, useEffect } from 'react';
import reactLogo from './assets/react.svg';
import './App.css';

function nvl(value1, value2) {
  if (value1 == null || value1.length == 0) return value2;
  return value1;
}

/*Root Component */
const App = () => {
  const [searchterm, setsearchterm] = useState(
    localStorage.getItem('search') || 'irgendwas'
  );

  //Daten aus der DB
  const [data, dataSetter] = useState([]);
  //Daten aus der DB laden mit fetch API
  const fetchData = async () => {
    const response = await fetch('http://localhost:8081/allitems')
      .then((response) => response.json())
      .then((data) => {
        dataSetter(data);
      });
  };
  //fetchData mit useEffect aufrufen
  useEffect(() => {
    fetchData();
  }, []);

  useEffect(() => {
    localStorage.setItem('search', searchterm);
  }, [searchterm]);

  const title = 'spaghetti';
  const arrOrig = [
    { key_id: 'alpha', value: 'butter' },
    { key_id: 'beta', value: 'keks' },
    { key_id: 'gammma', value: 'F1' },
  ];
  const arr = arrOrig.filter(function (item) {
    return item.key_id.includes(searchterm);
  });

  const dataFiltered = data.filter(function (item) {
    return item.Name.includes(searchterm);
  });

  return (
    <div>
      <h1>Vite + {title}</h1>
      <SearchTerm searchterm={searchterm} />
      <List list={arr} />
      <Search setter={setsearchterm} val={searchterm} />
      <DataList data={dataFiltered} setter={dataSetter} />
    </div>
  );
};

const DataList = ({ data, setter }) => {
  return (
    <div>
      <h2>data</h2>
      <ul>
        {data.map((x) => (
          <li>
            <i>{x.Vendor}</i>
            &nbsp;{x.Name}
          </li>
        ))}
      </ul>
    </div>
  );
};

const SearchTerm = ({ searchterm }) => <h2>{nvl(searchterm, '<leer>')}</h2>;

/* List Component */
const List = (props) => {
  return (
    <ol>
      {props.list.map((x) => (
        <Item key={x.key_id} item={x} />
      ))}
    </ol>
  );
};

/*Item Component*/
const Item = ({ item: { key_id, value } }) => (
  <li key={key_id}>
    <i>{key_id}</i>
    &nbsp;{value} hallo
  </li>
);

/* Search Component */
const Search = ({ setter, val }) => {
  const handleChange = (event) => {
    setter(event.target.value);
    console.log(event);
    console.log(event.target.value);
  };
  return (
    <div>
      <label htmlFor='search'>Search: </label>
      <input type='text' id='search' value={val} onChange={handleChange} />
    </div>
  );
};

export default App;
