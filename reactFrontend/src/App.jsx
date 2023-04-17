import { useState, useEffect } from 'react';
import reactLogo from './assets/react.svg';
import spinner from './assets/3.png';
import './App.css';

var baseurl = 'http://localhost:8081';
//var baseurl = '';

function nvl(value1, value2) {
  if (value1 == null || value1.length == 0) return value2;
  return value1;
}

function sleep(milliseconds) {
  const date = Date.now();
  let currentDate = null;
  do {
    currentDate = Date.now();
  } while (currentDate - date < milliseconds);
}

/*Root Component */
const App = () => {
  const [searchterm, setsearchterm] = useState(
    localStorage.getItem('search') || ''
  );
  const [isLoading, setIsLoading] = useState(false);

  //Daten aus der DB
  const [data, dataSetter] = useState([]);
  //Daten aus der DB laden mit fetch API
  const fetchData = async () => {
    setIsLoading(true);
    const response = await fetch(baseurl + '/allitems')
      .then((response) => {
        //sleep(3000);
        return response.json();
      })
      .then((data) => {
        console.log('waiting 3 sec');
        setTimeout(() => {
          dataSetter(data);
          console.log('done');
          setIsLoading(false);
        }, 3000);
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
      <h1>OrderHistory + {title}</h1>
      <SearchTerm searchterm={searchterm} />
      <List list={arr} />
      <Search setter={setsearchterm} val={searchterm} />
      <DataList data={dataFiltered} setter={dataSetter} load={isLoading} />
    </div>
  );
};

const DataList = ({ data, setter, load }) => {
  return (
    <div>
      <h2>data</h2>
      {load ? (
        <i>
          loading...
          <div id='fountainG'>
            <div id='fountainG_1' className='fountainG'></div>
            <div id='fountainG_2' className='fountainG'></div>
            <div id='fountainG_3' className='fountainG'></div>
            <div id='fountainG_4' className='fountainG'></div>
            <div id='fountainG_5' className='fountainG'></div>
            <div id='fountainG_6' className='fountainG'></div>
            <div id='fountainG_7' className='fountainG'></div>
            <div id='fountainG_8' className='fountainG'></div>
          </div>
        </i>
      ) : (
        <ul>
          {data.map((x) => (
            <li>
              <i>{x.Vendor}</i>
              &nbsp;{x.Name}
            </li>
          ))}
        </ul>
      )}
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
