import { useState, useEffect, useRef } from 'react';
import './App.css';
import ebaysvg from '/EBay_logo.svg';
import alternatesvg from '/Alternate.de_logo.svg';
import aliextrassvg from '/Aliexpress_logo.svg';

var baseurl = 'http://localhost:8081';
//var baseurl = '';

/*Root Component */
const App = () => {
  const getJsonFromLocalStorage = (key) => {
    const rawResult = localStorage.getItem(key);
    // console.log('from localStorage: %s -> %s', key, rawResult);
    if (rawResult && rawResult.length > 2) {
      // console.log('returning JSON: ' + rawResult);
      return JSON.parse(rawResult);
    }
    // console.log('returning null');
    return null;
  };

  const makeSet = (key) => {
    const json = getJsonFromLocalStorage(key);
    // console.log('makeSet: %s -> %s', key, json);
    // console.log('json null: ' + (json == null));
    // console.log('json undefined: ' + (json === undefined));
    if (json != null && json !== undefined) {
      // console.log('json length: ' + json.length);
    }
    if (json != null && json !== undefined && json.length >= 1) {
      // console.log('returning Set: ' + json);
      return new Set(json);
    } else {
      // console.log('returning empty Set');
      return new Set();
    }
  };

  const [searchterm, setsearchterm] = useState(
    localStorage.getItem('search') || ''
  );
  const [yearFilter, setYearFilter] = useState(makeSet('year'));
  const [vendorFilter, setVendorFilter] = useState(makeSet('vendor'));
  const [isLoading, setIsLoading] = useState(false);
  const [showForm, setShowForm] = useState(false);
  const modalAddItemDialog = useRef();

  const displayAddItemForm = () => {
    modalAddItemDialog.current.showModal();
  };

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
        }, 1000);
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
        dataSetter((oldData) => {
          //loop over data and update the data at id
          return oldData.map((item) => {
            if (item.Id === itemId) {
              return { ...item, ImgFile: data.imgFile };
            } else {
              return item;
            }
          });
        });
      })
      .catch((error) => {
        console.error('Error:', error);
      });
  };

  //Function to handle YearFilterChanges
  const handleYearFilterChange = (event) => {
    // console.log(event.target.id);
    const checkBoxes = document.querySelectorAll(
      'div.yearFilter input[type=checkbox]'
    );
    if (event.target.id == 0) {
      let allState = event.target.checked;
      for (var checkb of checkBoxes) {
        checkb.checked = allState;
      }
    }
    var currentYearsSet = new Set();
    for (var checkb of checkBoxes) {
      if (checkb.checked) {
        currentYearsSet.add(Number(checkb.id));
      } else {
        currentYearsSet.delete(Number(checkb.id));
      }
    }
    setYearFilter((yearFilter) => (yearFilter = currentYearsSet));
  };

  //Function to handle VendorFilterChanges
  const handleVendorFilterChange = (event) => {
    // console.log(event.target.id);
    const checkBoxes = document.querySelectorAll(
      'div.vendorFilter input[type=checkbox]'
    );
    if (event.target.id == 'All') {
      let allState = event.target.checked;
      for (var checkb of checkBoxes) {
        checkb.checked = allState;
      }
    }
    var currentVendorsSet = new Set();
    for (var checkb of checkBoxes) {
      if (checkb.checked) {
        currentVendorsSet.add(checkb.id);
      } else {
        currentVendorsSet.delete(checkb.id);
      }
    }
    setVendorFilter((vendorFilter) => (vendorFilter = currentVendorsSet));
  };

  useEffect(() => {
    localStorage.setItem('search', searchterm);
  }, [searchterm]);
  useEffect(() => {
    localStorage.setItem('year', JSON.stringify([...yearFilter]));
  }, [yearFilter]);
  useEffect(() => {
    localStorage.setItem('vendor', JSON.stringify([...vendorFilter]));
  }, [vendorFilter]);

  const dataFiltered = data.filter(function (item) {
    return item.Name.toLowerCase().includes(searchterm.toLowerCase()) &&
      (yearFilter.has(new Date(item.PurchaseDate).getFullYear()) ||
        (yearFilter.size == 0 ? true : false)) &&
      (vendorFilter.has(item.Vendor) || vendorFilter.size == 0)
      ? true
      : false;
  });

  return (
    <div>
      <h1>OrderHistory</h1>
      <button onClick={displayAddItemForm}>Add Item</button>
      <OverlayForm show={showForm} modalDialog={modalAddItemDialog} />
      <Search setter={setsearchterm} val={searchterm} />
      <YearFilter
        data={data}
        yearFilter={yearFilter}
        handleYearFilterChange={handleYearFilterChange}
      />
      <VendorFilter
        data={data}
        vendorFilter={vendorFilter}
        handleVendorFilterChange={handleVendorFilterChange}
      />
      <DataList data={dataFiltered} load={isLoading} handleFile={handleFile} />
    </div>
  );
};

const AddItemButton = ({ handleShow }) => (
  <button onClick={handleShow}>Add Item</button>
);

//----------------------------------------------------------------------

//Compoent: Overlay Form for adding new Items
const OverlayForm = ({ show, modalDialog }) => {
  const handleClose = (event) => {
    event.preventDefault();
    modalDialog.current.close();
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    const form = event.target;
    const formData = new FormData(form);
    fetch(baseurl + '/newItemManual', {
      method: 'POST',
      body: formData,
    })
      .then((response) => response.json())
      .then((data) => {
        console.log('Returnmessage:', data);
        handleAdd(data);
        handleClose();
      })
      .catch((error) => {
        console.error('Error:', error);
      });
  };

  return (
    <dialog data-modal className='modal dialog' ref={modalDialog}>
      <form method='post' onSubmit={handleSubmit} className='itemDialog'>
        <h2 className='form-header'>Add new Item</h2>
        <label className='form-itemLabel' htmlFor='itemName'>
          Name
        </label>
        <input className='form-item' type='text' name='itemName' id='name' />

        <label className='form-priceLabel' htmlFor='price'>
          Price
        </label>
        <input
          className='form-price'
          type='text'
          name='price'
          id='price'
          placeholder='3,99'
        />

        <label className='form-dateLabel' htmlFor='date'>
          Date
        </label>
        <input
          className='form-date'
          type='text'
          name='date'
          id='date'
          placeholder='31.12.2001'
        />

        <label className='form-platformLabel' htmlFor='platform'>
          Platform
        </label>
        <input
          className='form-platform'
          type='platform'
          name='platform'
          id='platform'
        />

        <label className='form-currencyLabel' htmlFor='currency'>
          Currency
        </label>
        <input
          className='form-currency'
          type='text'
          name='currency'
          id='currency'
          placeholder='€'
        />

        <label className='form-imgUrlLabel' htmlFor='imgUrl'>
          Image Url
        </label>
        <input
          className='form-imgUrl'
          type='text'
          name='imgUrl'
          id='imgUrl'
          placeholder='https://example.com/test.png'
        />

        <label className='form-divLabel' htmlFor='div'>
          Div
        </label>
        <input
          className='form-div'
          type='text'
          name='div'
          id='div'
          placeholder='{"key":"value", "key2":"value2"}'
        />

        <span className='form-button'>
          <button data-close-modal onClick={handleClose}>
            Close
          </button>
          <button type='submit'>Submit</button>
        </span>
      </form>
    </dialog>
  );
};

//Component: VendorFilter
const VendorFilter = ({ data, vendorFilter, handleVendorFilterChange }) => {
  var uniqueVendors = [...new Set(data.map((item) => item.Vendor))];
  uniqueVendors.unshift('All');
  return (
    <div className='vendorFilter'>
      {uniqueVendors.map((vendor) => (
        <div key={vendor}>
          <input
            type='checkbox'
            value={vendor}
            id={vendor}
            onChange={handleVendorFilterChange}
            checked={vendorFilter.has(vendor)}
          />
          <label htmlFor={vendor}>{vendor}</label>
        </div>
      ))}
    </div>
  );
};

//Component: YearFilter
const YearFilter = ({ data, yearFilter, handleYearFilterChange }) => {
  var uniqueYears = [
    ...new Set(data.map((item) => new Date(item.PurchaseDate).getFullYear())),
  ];
  uniqueYears.unshift(0);
  return (
    <div className='yearFilter'>
      {uniqueYears.map((year) => (
        <div key={year}>
          <input
            type='checkbox'
            value={year === 0 ? 'All' : year}
            id={year}
            onChange={handleYearFilterChange}
            checked={yearFilter.has(year)}
          />
          <label htmlFor={year}>{year === 0 ? 'All' : year}</label>
        </div>
      ))}
    </div>
  );
};

const DataList = ({ data, load, handleFile }) => {
  return (
    <div>
      <h2>{data.length} Orders:</h2>
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
  const handleClick = (event) => {
    hiddenFileInput.current.click();
  };
  const handleChange = (event) => {
    const fileUploaded = event.target.files[0];
    if (event.target.files[0] !== undefined) {
      handleFile(fileUploaded, itemId, vendor);
    }
  };
  return (
    <>
      <img
        className='picture'
        src={`${baseurl}/img/${ImgFile}`}
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
