// ==UserScript==
// @name         ebay purchases client script
// @namespace    http://tampermonkey.net/
// @version      0.1
// @description  ebay purchases client script
// @author       You
// @match        https://www.ebay.de/mye/myebay/*
// @icon         data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==
// @grant        none
// ==/UserScript==

async function postData(url = '', data = {}) {
  // Default options are marked with *
  const response = await fetch(url, {
    method: 'POST', // *GET, POST, PUT, DELETE, etc.
    mode: 'cors', // no-cors, *cors, same-origin
    cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
    credentials: 'same-origin', // include, *same-origin, omit
    headers: {
      'Content-Type': 'application/json',
      // 'Content-Type': 'application/x-www-form-urlencoded',
    },
    redirect: 'follow', // manual, *follow, error
    referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
    body: JSON.stringify(data), // body data type must match "Content-Type" header
  });
  return response.json(); // parses JSON response into native JavaScript objects
}

(function () {
  'use strict';
  let orderList = document.querySelectorAll(
    '#container-1-anchor > section > div.m-container__body > div.m-container__body--items > div > div > div.m-ph-card__content.m-ph-layout-row-hidable > div'
  );
  console.log('Items auf Seite: %s', orderList.length);
  for (var order of orderList) {
    console.log('working on order...');
    //artikelnummer:
    // let artikelnummer = order.querySelector('[data-listing-id]');
    // console.log(artikelnummer.getAttribute('data-listing-id'));
    //itemName:
    let itemNameElement = order.querySelector(
      'div.m-container-item-layout-row.m-container-item-layout-row__body > div.m-container-item-layout-row__wrapper > div.container-item-col.container-item-col-item-info > div.container-item-col__info-item-info-primary.container-item-col__info-item-info-title > div > a'
    );
    console.log('itemName: %o', itemNameElement.text);
    //price
    let itemPrice = order.querySelector(
      'div.m-container-item-layout-row.m-container-item-layout-row__body > div.m-container-item-layout-row__wrapper > div.container-item-col.container-item-col-item-info > div.container-item-col__info-item-info-primary.container-item-col__info-item-info-additionalPrice > div > span'
    );
    let priceNumber = /^[^\d]*(\d.+)/.exec(itemPrice.innerHTML);
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
    console.log('currency: %d', itemCurrency[0]);
    let currency;
    if (itemCurrency[0]) {
      currency = itemCurrency[0].replace('EUR', '€').replace('US', '$');
    }
    //vendor
    let vendorElement = order.querySelector(
      'div.m-container-item-layout-row.m-container-item-layout-row__body > div.m-container-item-layout-row__wrapper > div.container-item-col.container-item-col-item-info > div.container-item-col__info-item-info-primary.container-item-col__info-item-info-sellerInfo > div:nth-child(2) > a > span'
    );
    console.log('Haendler: %s', vendorElement.firstChild.textContent);
    //imgUrl
    let imgElement = order.querySelector('div.m-image img[src][alt]');
    console.log('SRC: ' + imgElement.getAttribute('src'));
    let imgUrl = imgElement.getAttribute('src');
    if (imgElement.getAttribute('data-imgurl')) {
      console.log('DATA-IMGURL: ' + imgElement.getAttribute('data-imgurl'));
      imgUrl = imgElement.getAttribute('data-imgurl');
    }
    //orderDate
    let dateElement = order.querySelector('div.ph-col__info-orderDate dd');
    console.log(dateElement.textContent);
    //build object for later json marshal
    var orderObj = {
      artikelnummer: parseInt(artikelnummer.getAttribute('data-listing-id')),
      itemName: itemNameElement.text,
      price: parseInt(price),
      vendor: vendorElement.firstChild.textContent,
      imgUrl: imgUrl,
      purchaseDate: dateElement.textContent
        .replace('Mai', 'May')
        .replace('Okt', 'Oct')
        .replace('Dez', 'Dec')
        .replace('Mär', 'Mar'),
      currency: currency,
    };
    console.log(JSON.stringify(orderObj));
    //request to api

    // postData('http://localhost:8081/order/ebay', orderObj)
    //   .then(data => {
    //     console.log(data); // JSON data parsed by `data.json()` call
    //   });
  }
})();
