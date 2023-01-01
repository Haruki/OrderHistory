// ==UserScript==
// @name         aliexpress purchases client script
// @namespace    http://tampermonkey.net/
// @version      0.1
// @description  aliexpress purchases client script
// @author       You
// @match        https://www.aliexpress.com/p/order/*
// @grant        none
// @run-at       document-idle
// ==/UserScript==

async function postData(url = '', data = {}) {
    // Default options are marked with *
    const response = await fetch(url, {
        method: 'POST', // *GET, POST, PUT, DELETE, etc.
        mode: 'cors', // no-cors, *cors, same-origin
        cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
        credentials: 'same-origin', // include, *same-origin, omit
        headers: {
            'Content-Type': 'application/json'
            // 'Content-Type': 'application/x-www-form-urlencoded',
        },
        redirect: 'follow', // manual, *follow, error
        referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
        body: JSON.stringify(data) // body data type must match "Content-Type" header
    });
    return response.json(); // parses JSON response into native JavaScript objects
}


//Achtung aliexpress sonderlösung
//bei ebay und alternate läuft der upload der orders automatisch
//bei aliexpress muss nach dem vollständigen aufklappen aller orders der process über den button manuell gestartet werden, da sonst beim automatischen lauf nicht alle orders sichtbar sind
//der Button wird oben links unter "Konto" eingeblendet nur mit dem Text "process Orders"
const createButton = function() {
    console.log("start createButton")
    let navTitleElement = document.querySelector("#page-menu div.page-menu-item.page-menu-title")
    const newButton = document.createElement("button");
    newButton.setAttribute("type","button")
    newButton.addEventListener("click", processOrders)
    newButton.innerText = "process Orders"
    //after: Siehe https://stackoverflow.com/questions/4793604/how-to-insert-an-element-after-another-element-in-javascript-without-using-a-lib
    // runterscrollen zu 2022 solution
    navTitleElement.after(newButton)
}

window.addEventListener('load', createButton, false)

const processOrders = function() {
    'use strict';
    console.log("test")
    let orderList = document.querySelectorAll("div.order-item")
    console.log("orderListLength: " + orderList.length)
    for (var order of orderList) {
        //bestellnummer:
        console.log("working on order...")
        let bestellnummerElem = order.querySelector("div.order-item-header div.order-item-header-right div div:nth-child(2)")
        let bestellnummerRegex = /^[^\d]*(\d{14,16})/.exec(bestellnummerElem.textContent)
        let bestellnummer
        if (bestellnummerRegex[1]) {
          bestellnummer = bestellnummerRegex[1]
        } else {
          bestellnummer = "unbekannt"
        }
        console.log("bestellnummer: " + bestellnummer)
        //itemName:
        let itemElement = order.querySelector("div.order-item-content div.order-item-content-body div div.order-item-content-info-name a span")
        let itemName
        if (itemElement && itemElement.getAttribute("title")) {
          console.log(itemElement.getAttribute("title"))
          itemName = itemElement.getAttribute("title")
        } else {
          itemName = "unbekannte Bezeichnung"
        }
        //itemOption
        let itemOptionElem = order.querySelector("div.order-item-content div.order-item-content-body div div.order-item-content-info-sku")
        let itemOption
        if (itemOptionElem) {
            itemOption = itemOptionElem.textContent
            console.log("itemOption: " + itemOption)
        }
        //Einzelprice
        let singlePriceElem = order.querySelector("div.order-item-content div.order-item-content-body div div.order-item-content-info-number")
        let singlePrice
        if (singlePriceElem) {
            let singlePriceAr = /^[^\d]*(\d+[\,|\.]\d+)?.*$/.exec(singlePriceElem.textContent)
            if (singlePriceAr[1]) {
              singlePrice = singlePriceAr[1].replace(",", "").replace(".","")
              console.log("singlePrice: " + singlePrice)
            }
        }
        //count
        let countElem = order.querySelector("div.order-item-content div.order-item-content-body div div.order-item-content-info-number span")
        let count
        if (countElem) {
            let countAr = /^[^\d]*([\d]*)$/.exec(countElem.textContent)
            if (countAr[1]) {
              count = countAr[1]
              console.log("count: " + count)
            }
        }
        //price
        let itemPrice = order.querySelector("div.order-item-content div.order-item-content-opt div.order-item-content-opt-price span")
        let priceNumber = /^[^\d]*(\d+[\,|\.]\d+)?$/.exec(itemPrice.textContent)
        let price
        if (priceNumber[1]) {
            console.log("regex price raw: " + typeof priceNumber[1])
            price = priceNumber[1].replace(",", "").replace(".","")
        } else {
            console.log("non regex price: " + itemPrice.innerHTML.substring(4, ).replace(",", ""))
            price = itemPrice.innerHTML.substring(4, ).replace(",", "")
        }
        console.log("price: " + price)
        //currency
        let itemCurrency = /(i?)([$£€¥])/.exec(itemPrice.innerHTML)
        console.log("currency: " + itemCurrency[0])
        let currency
        if (itemCurrency[0]) {
            currency = itemCurrency[0].replace("EUR","€").replace("US","$")
        }
        //vendor
        let vendorElement = order.querySelector("div.order-item-store span a span")
        let vendor = vendorElement.textContent
        console.log("Vendor:" + vendor)
        //imgUrl
        let imgElement = order.querySelector("div.order-item-content div.order-item-content-body a div")
        let imgFullText = imgElement.getAttribute("style")
        let imgUrlAr = /https?:\/\/[^\s"]+/.exec(imgFullText)
        console.log("imgUrlAr: " + imgUrlAr)
        let imgUrl
        if (imgUrlAr[1]) {
            imgUrl = encodeURIComponent(imgUrlAr[1])
            console.log("imgUrl: " + imgUrl)
        } else {
            imgUrl = encodeURIComponent(imgUrlAr)
        }
        //orderDate
        let orderDate
        let dateElement = order.querySelector("div.order-item-header div.order-item-header-right div div:nth-child(1)")
        let orderDateRegex = /[^\d]*(\d+\.\s[\w|ä]+\s\d+).*$/.exec(dateElement.textContent) //achtung: Ergebnis ist 1+ capture groups als array
        if (orderDateRegex[1]) {
            orderDate = orderDateRegex[1]
        } else {
            orderDate = orderDateRegex
            console.log("non array orderDate: " + typeof orderDate)
        }
        orderDate = orderDate.replace("Mai", "May").replace("Okt","Oct").replace("Dez","Dec").replace("Mär","Mar")
        if (orderDate.length == 11) {
            orderDate = "0" + orderDate
        }
        console.log("Bestelldatum: " + orderDate)
        //build object for later json marshal
        var orderObj = {
            bestellnummer: parseInt(bestellnummer),
            itemName: itemName,
            itemOption: itemOption,
            price: parseInt(price),
            vendor: vendor,
            imgUrl: imgUrl,
            purchaseDate: orderDate,
            currency: currency,
            singlePrice: parseInt(singlePrice),
            anzahl: parseInt(count)
        }
        console.log(JSON.stringify(orderObj))
        //request to api

                postData('http://localhost:8081/order/aliexpress', orderObj)
          .then(data => {
            console.log(data); // JSON data parsed by `data.json()` call
          });
    }

}