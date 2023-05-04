import { useState, useEffect, useRef } from 'react';
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

  //Function for fileupload to api
  const handleFile = async (file, itemId, vendor) => {
    let formData = new FormData();
    formData.append('file', file);
    formData.append('itemId', itemId);
    formData.append('vendor', vendor);
    fetch(baseurl + '/imageUpload', {
      method: 'POST',
      body: formData,
    })
      .then((response) => response.json())
      .then((data) => {
        console.log('Success:', data);
      })
      .catch((error) => {
        console.error('Error:', error);
      });
  };

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
      <DataList data={dataFiltered} load={isLoading} handleFile={handleFile} />
    </div>
  );
};

const DataList = ({ data, load, handleFile }) => {
  return (
    <div>
      <h2>Orders:</h2>
      {load ? (
        <LoadingHint />
      ) : (
        <div className='datalist'>
          {data.map((x) => (
            <OrderItem key={x.Id} {...x} handleFile={handleFile} />
          ))}
        </div>
      )}
    </div>
  );
};

const OrderItem = ({
  Id,
  Vendor,
  Name,
  PurchaseDate,
  Price,
  Currency,
  ImgFile,
  Div,
  handleFile,
  itemId,
}) => {
  return (
    <article className='entry orderItem'>
      <Picture
        ImgFile={ImgFile}
        handleFile={handleFile}
        itemId={Id}
        vendor={Vendor}
      />
      <span className='entry artikel'>{Name}</span>
      {/*<span className='entry platform'>{Vendor}</span>*/}
      <Platform Vendor={Vendor} />
      <span className='entry purchaseDate'>{PurchaseDate}</span>
      <span className='entry price'>{Price / 100}</span>
      <span className='entry currency'>{Currency}</span>
      <Sonstiges Div={Div} />
    </article>
  );
};

const Sonstiges = ({ Div }) => {
  var divObject = JSON.parse(Div);
  var divEntries = Object.entries(divObject);
  return (
    <div className='entry sonstiges'>
      {divEntries.map(([key, value]) => (
        <span key={key}>
          {key}:{' '}
          {key.toLowerCase().indexOf('preis') == -1 ? value : value / 100}
        </span>
      ))}
    </div>
  );
};

const Picture = ({ ImgFile, handleFile, itemId, vendor }) => {
  const hiddenFileInput = useRef(null);
  const getFileName = (ImgFile) => {
    let lastSlashIndex = ImgFile.lastIndexOf('/');
    let filename = ImgFile.substring(lastSlashIndex + 1, ImgFile.length);
    return filename;
  };
  const handleClick = (event) => {
    hiddenFileInput.current.click();
  };
  const handleChange = (event) => {
    const fileUploaded = event.target.files[0];
    handleFile(fileUploaded, itemId, vendor);
  };
  return (
    <>
      <img
        className='picture'
        src={`${baseurl}/img/backup/${getFileName(ImgFile)}`}
        onClick={handleClick}
      />
      <input
        type='file'
        ref={hiddenFileInput}
        onChange={handleChange}
        style={{ display: 'none' }}
      />
    </>
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
