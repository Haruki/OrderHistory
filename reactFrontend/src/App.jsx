import { useState, useEffect } from 'react';
import './App.css';

var baseurl = 'http://localhost:8081';
//var baseurl = '';

function nvl(value1, value2) {
  if (value1 == null || value1.length == 0) return value2;
  return value1;
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

  const dataFiltered = data.filter(function (item) {
    return item.Name.includes(searchterm);
  });

  return (
    <div>
      <h1>OrderHistory</h1>
      <SearchTerm searchterm={searchterm} />
      <Search setter={setsearchterm} val={searchterm} />
      <DataList data={dataFiltered} load={isLoading} />
    </div>
  );
};

const DataList = ({ data, load }) => {
  return (
    <div>
      <h2>Orders:</h2>
      {load ? (
        <LoadingHint />
      ) : (
        <ul>
          {data.map((x) => (
            <li key={x.Id}>
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

const LoadingHint = () => (
  <>
    <i>loading...</i>
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
  </>
);

export default App;
