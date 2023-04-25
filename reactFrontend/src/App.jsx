import { useState, useEffect } from 'react';
import './App.css';
import ebaysvg from '/EBay_logo.svg';
import alternatesvg from '/Alternate.de_logo.svg';
import aliextrassvg from '/Aliexpress_logo.svg';

var baseurl = 'http://localhost:8081';
//var baseurl = '';

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
    return item.Name.toLowerCase().includes(searchterm.toLowerCase());
  });

  return (
    <div>
      <h1>OrderHistory</h1>
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
        <div className='datalist'>
          {data.map((x) => (
            <OrderItem key={x.Id} {...x} />
          ))}
        </div>
      )}
    </div>
  );
};

const OrderItem = ({
  Vendor,
  Name,
  PurchaseDate,
  Price,
  Anzahl,
  Currency,
  ImgFile,
}) => {
  return (
    <article className='entry orderItem'>
      <Picture ImgFile={ImgFile} />
      <span className='entry artikel'>{Name}</span>
      {/*<span className='entry platform'>{Vendor}</span>*/}
      <Platform Vendor={Vendor} />
      <span className='entry purchaseDate'>{PurchaseDate}</span>
      <span className='entry anzahl'>{Anzahl}</span>
      <span className='entry price'>{Price / 100}</span>
      <span className='entry currency'>{Currency}</span>
      <span className='entry sonstiges'>lalala sonstiges lalal</span>
    </article>
  );
};

const Picture = ({ ImgFile }) => {
  const getFileName = (ImgFile) => {
    let lastSlashIndex = ImgFile.lastIndexOf('/');
    let filename = ImgFile.substring(lastSlashIndex + 1, ImgFile.length);
    return filename;
  };
  return (
    <img
      className='picture'
      src={`${baseurl}/img/backup/${getFileName(ImgFile)}`}
    />
  );
};

const Platform = ({ Vendor }) => {
  switch (Vendor) {
    case 'ebay':
      return <img className='entry platform' src={ebaysvg} alt='ebay' />;
    case 'alternate':
      return (
        <img className='entry platform' src={alternatesvg} alt='alternate' />
      );
    case 'aliexpress':
      return (
        <img className='entry platform' src={aliextrassvg} alt='aliexpress' />
      );
    default:
      return <span className='entry platform'>{Vendor}</span>;
  }
};

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
