// ==UserScript==
// @name         amazon purchases client script
// @namespace    http://tampermonkey.net/
// @version      0.1
// @description  amazon purchases client script
// @author       You
// @match        https://www.amazon.de/gp/your-account/order-details*
// @icon         data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==
// @require     file://D:\data\coding\OrderHistory\client\tampermonkey\amazon_purchases_client.js
// @grant        none
// @run-at       document-idle
// ==/UserScript==

var baseUrl = 'http://localhost:8081';

function buildButton(parent, data) {
  parent.innerHTML = '';
  var button = document.createElement('button');
  button.innerHTML = 'Upload data';
  button.disabled = true;
  button.addEventListener('click', function () {
    fetchData(baseUrl + '/newItemManual', data, 'POST').then((response) => {
      console.log('order/amazon: %s', response.message); // JSON data parsed by `data.json()` call
      if (response.message == 'Item added') {
        button.innerHTML = 'Upload successful';
        button.disabled = true;
        button.style.backgroundColor = 'green';
      } else {
        button.innerHTML = 'Upload failed';
        button.disabled = true;
        button.style.backgroundColor = 'red';
      }
    });
  });
  parent.appendChild(button);
  //fetch current state from api
  fetch(
    baseUrl +
      '/checkItemExists?' +
      new URLSearchParams({
        itemName: data.itemName,
        purchaseDate: data.date,
        vendor: 'amazon',
      })
  ).then((response) => {
    console.log(response.status);
    response.json().then((data) => {
      console.log(data);
      if (response.status == 200) {
        button.innerHTML = 'Item exists';
        button.disabled = true;
        button.style.backgroundColor = 'green';
      } else {
        button.innerHTML = 'Upload data';
        button.disabled = false;
        button.style.backgroundColor = 'red';
      }
    });
  });
}

async function fetchData(url = '', data = {}, method = 'POST') {
  // Default options are marked with *
  const formData = new FormData();
  for (const key in data) {
    formData.append(key, data[key]);
  }
  const response = await fetch(url, {
    method: method, // *GET, POST, PUT, DELETE, etc.
    mode: 'cors', // no-cors, *cors, same-origin
    cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
    credentials: 'same-origin', // include, *same-origin, omit
    redirect: 'follow', // manual, *follow, error
    referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
    body: formData, // body data type must match "Content-Type" header
  });
  return response.json(); // parses JSON response into native JavaScript objects
}

function convertDate(dateString) {
  dateString = dateString.replace('Bestellt am ', '');
  dateString = dateString.replace('.', '');
  const months = {
    Januar: '01',
    Februar: '02',
    März: '03',
    April: '04',
    Mai: '05',
    Juni: '06',
    Juli: '07',
    August: '08',
    September: '09',
    Oktober: '10',
    November: '11',
    Dezember: '12',
  };
  const [day, month, year] = dateString.split(' ');
  const formattedDay = day.padStart(2, '0');
  const formattedMonth = months[month];
  return `${formattedDay}.${formattedMonth}.${year}`;
}

(function () {
  'use strict';
  //global order date in case the per item date is not available
  let orderDate = document.querySelector('.order-date-invoice-item').innerText;
  if (orderDate) {
    orderDate = convertDate(orderDate);
  }
  console.log('global orderDate: %s', orderDate);
  //get all orders -> List of orders
  let orderList = document.querySelectorAll(
    '.shipment div.a-fixed-left-grid.a-spacing-base,.shipment div.a-fixed-left-grid.a-spacing-none'
  );
  console.log('Items auf Seite: %s', orderList.length);
  for (var order of orderList) {
    console.log('working on order...');
    //orderDate
    let dateElement = order.querySelector(
      'span.a-size-medium.a-color-base.a-text-bold'
    );
    let dateCleaned;
    if (dateElement) {
      console.log('orderDate: %s', dateElement.textContent);
      dateCleaned = dateElement.textContent.replace('Zugestellt:', '').trim();
    } else {
      dateCleaned = orderDate;
    }
    //itemName:
    // let itemNameElement = order.querySelector(
    //   'a.a-link-normal[href^="/gp/product"]'
    // );
    let itemName;
    let itemNameElement = order.querySelectorAll(
      // 'div > div.a-fixed-right-grid.a-spacing-top-medium > div > div.a-fixed-right-grid-col.a-col-left > div > div > div > div.a-fixed-left-grid-col.yohtmlc-item.a-col-right > div:nth-child(1) > a'
      'a[href^="/gp/product/"]'
    )[1];
    itemName = itemNameElement.innerText.trim();
    console.log('itemName: ' + itemName);
    //price

    let itemPrice = order.querySelector('span.a-color-price nobr');
    let priceNumber = /^[^\d]*(\d.+)/.exec(itemPrice.innerText.trim());
    let price;
    if (priceNumber[1]) {
      console.log('regex price raw: ' + typeof priceNumber[1]);
      console.log('regex price number: ' + priceNumber[1].replace(',', ''));
      price = priceNumber[1].replace(',', '');
    } else {
      console.log(
        'non regex price: ' + itemPrice.innerHTML.substring(4).replace(',', '')
      );
      price = itemPrice.innerHTML.substring(4).replace(',', '');
    }

    //currency
    let itemCurrency = /^[^(\s|\d)]*/.exec(itemPrice.innerHTML);
    console.log('currency: %s', itemCurrency[0]);
    let currency;
    if (itemCurrency[0]) {
      currency = itemCurrency[0].replace('EUR', '€').replace('US', '$');
    }
    //vendor
    let vendor;
    let vendorElement = document
      .evaluate(
        '//span[contains(., "Verkauf durch")]',
        order,
        null,
        XPathResult.ANY_TYPE,
        null
      )
      .iterateNext();
    // let vendorElement = order.querySelector(
    //   'div > div.a-fixed-right-grid.a-spacing-top-medium > div > div.a-fixed-right-grid-col.a-col-left > div > div > div > div.a-fixed-left-grid-col.yohtmlc-item.a-col-right > div:nth-child(2) > span > a'
    // );
    if (vendorElement) {
      console.log('Haendler: %s', vendorElement.firstChild.textContent.trim);
      vendor = vendorElement.firstChild.textContent
        .replace('Verkauf durch:', '')
        .trim();
    }
    //imgUrl
    let imgElement = order.querySelector(
      // 'div > div.a-fixed-right-grid.a-spacing-top-medium > div > div.a-fixed-right-grid-col.a-col-left > div > div > div > div.a-text-center.a-fixed-left-grid-col.a-col-left > div > a > img'
      'img[alt=""]'
    );
    console.log('SRC: ' + imgElement.getAttribute('src'));
    let imgUrl = imgElement.getAttribute('src');
    if (imgElement.getAttribute('data-imgurl')) {
      console.log('DATA-IMGURL: ' + imgElement.getAttribute('data-imgurl'));
      imgUrl = imgElement.getAttribute('data-imgurl');
    }

    //build object for later json marshal
    //build div object
    let divObj = {};
    if (vendor) {
      divObj.vendor = vendor;
    }
    var orderObj = {
      //   artikelnummer: parseint(artikelnummer.getattribute('data-listing-id')),
      itemName: itemName,
      price: parseInt(price),
      imgUrl: imgUrl,
      date: dateCleaned,
      platform: 'amazon',
      currency: currency,
      div: JSON.stringify(divObj),
    };
    console.log(JSON.stringify(orderObj));

    //build button
    var parent = order.querySelector(
      // 'div > div.a-fixed-right-grid.a-spacing-top-medium > div > div.a-fixed-right-grid-col.a-col-left > div > div > div > div.a-fixed-left-grid-col.yohtmlc-item.a-col-right > div:nth-child(6) > span'
      'span[data-action="bia_button"], div.a-row span.a-size-small div.a-row.a-size-small, div.a-row span.a-declarative[data-action="a-popover"]'
    );
    buildButton(parent, orderObj);
  }
})();
