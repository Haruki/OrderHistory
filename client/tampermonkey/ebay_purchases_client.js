// ==UserScript==
// @name         ebay purchases client script
// @namespace    http://tampermonkey.net/
// @version      0.1
// @description  ebay purchases client script
// @author       You
// @match        https://www.ebay.de/mye/myebay/*
// @icon         data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==
// @require     file://D:\data\coding\OrderHistory\client\tampermonkey\ebay_purchases_client.js
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
    fetchData(baseUrl + '/order/ebay', data, 'POST').then((response) => {
      console.log('order/ebay: %s', response.message); // JSON data parsed by `data.json()` call
      if (response.message == 'Success') {
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
        purchaseDate: data.purchaseDate,
        vendor: 'ebay',
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

async function fetchData(url = '', data = {}, method = 'GET') {
  // Default options are marked with *
  const response = await fetch(url, {
    method: method, // *GET, POST, PUT, DELETE, etc.
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
  let orderList = document.querySelectorAll('.m-ph-card.m-order-card');
  console.log('Items auf Seite: %s', orderList.length);
  for (var order of orderList) {
    console.log('working on order...');
    //orderDate
    let dateElement = order.querySelector(
      'div.primary__item.secondaryMessage span:nth-child(2).primary__item--item-text'
    );
    console.log('orderDate: %s', dateElement.textContent);
    let dateCleaned = dateElement.textContent
      .replace('Mai', 'May')
      .replace('Okt', 'Oct')
      .replace('Dez', 'Dec')
      .replace('Mär', 'Mar');
    //loop item cards
    let itemCards = order.querySelectorAll(
      '.m-item-card.m-container-item-layout-row__body'
    );
    console.log('Items in order: %s', itemCards.length);
    for (var itemCard of itemCards) {
      //itemName:
      let itemNameElement = itemCard.querySelector('div a.nav-link');
      console.log('itemName: %o', itemNameElement.text);
      //price
      let itemPrice = itemCard.querySelector(
        'div.container-item-col__info-item-info-additionalPrice div span'
      );
      let priceNumber = /^[^\d]*(\d.+)/.exec(itemPrice.innerHTML);
      let price;
      if (priceNumber[1]) {
        console.log('regex price raw: ' + typeof priceNumber[1]);
        console.log('regex price number: ' + priceNumber[1].replace(',', ''));
        price = priceNumber[1].replace(',', '');
      } else {
        console.log(
          'non regex price: ' +
            itemPrice.innerHTML.substring(4).replace(',', '')
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
      let vendorElement = itemCard.querySelector(
        'div.container-item-col__info-item-info-primary.container-item-col__info-item-info-sellerInfo div a span'
      );
      console.log('Haendler: %s', vendorElement.firstChild.textContent);
      //imgUrl
      let imgElement = itemCard.querySelector('div.m-image img[src][alt]');
      console.log('SRC: ' + imgElement.getAttribute('src'));
      let imgUrl = imgElement.getAttribute('src');
      if (imgElement.getAttribute('data-imgurl')) {
        console.log('DATA-IMGURL: ' + imgElement.getAttribute('data-imgurl'));
        imgUrl = imgElement.getAttribute('data-imgurl');
      }
      //build object for later json marshal
      var orderObj = {
        //   artikelnummer: parseint(artikelnummer.getattribute('data-listing-id')),
        itemName: itemNameElement.text,
        price: parseInt(price),
        ebaySpecial: { Haendler: vendorElement.firstChild.textContent },
        imgUrl: imgUrl,
        purchaseDate: dateCleaned,
        currency: currency,
      };
      console.log(JSON.stringify(orderObj));

      //build button
      var parent = itemCard.querySelector('.fake-menu-button');
      buildButton(parent, orderObj);
    }
  }
})();
