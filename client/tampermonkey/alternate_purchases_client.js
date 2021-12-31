// ==UserScript==
// @name         alternate purchases client script
// @namespace    http://tampermonkey.net/
// @version      0.1
// @description  alternate purchases client script
// @author       You
// @match        https://www.alternate.de/*
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
            'Content-Type': 'application/json'
            // 'Content-Type': 'application/x-www-form-urlencoded',
        },
        redirect: 'follow', // manual, *follow, error
        referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
        body: JSON.stringify(data) // body data type must match "Content-Type" header
    });
    return response.json(); // parses JSON response into native JavaScript objects
}

window.addEventListener('load', function() {
    'use strict';
    console.log("test")
    let orderList = document.querySelectorAll("#order-details > div > div.col-12.col-lg-10 > div:nth-child(1) > div.col-12.mt-5 > div.d-table.w-100 > a")
    console.log(orderList.length)
    //orderDate (global for all items -> "Detailansicht")
    let dateElement = document.querySelector("#order-details > div > div.col-12.col-lg-10 > div:nth-child(1) > div.col-12.col-md-9 > div:nth-child(1) > div:nth-child(2) > div > span:nth-child(1)")
    let purchaseDate = dateElement.textContent
    console.log("Purchase Date: " + purchaseDate)

    for (var order of orderList) {
        //artikelnummer:
        console.log("working on order...")
        /*
        let artikelnummer = order.querySelector("[data-listing-id]")
        console.log(artikelnummer.getAttribute("data-listing-id"))
        */
        //itemName:
        let itemElement = order.querySelector("div.d-table-cell > div > div.col-10 > div > div.col-12.font-weight-bold")
        let itemName = itemElement.textContent.trim()
        console.log(itemName)
        //price
        //let itemPrice = order.querySelector("div:nth-child(3) > span")
        //console.log(order.innerHTML)
        let itemPrice = order.querySelector("div:nth-child(2) > span")
        console.log("price raw: " + itemPrice.textContent)
        let price = 0
        if (itemPrice.textContent.length) {
            let priceNumber = /^[^\d]*(\d.+)/.exec(itemPrice.textContent)
            if (priceNumber[1]) {
                console.log("regex price raw: " + typeof priceNumber[1])
                console.log("regex price number: " + priceNumber[1].replace(",", ""))
                price = priceNumber[1].replace(",", "")
            } else {
                console.log("non regex price: " + itemPrice.innerHTML.substring(4, ).replace(",", ""))
                price = itemPrice.innerHTML.substring(4, ).replace(",", "")
            }
        }

        //currency
        if (price > 0) {
            let itemCurrency = /^[^(\s|\d)]*/.exec(itemPrice.innerHTML)
            console.log("currency: " + itemCurrency[0])
            let currency
            if (itemCurrency[0]) {
                currency = itemCurrency[0].replace("EUR","€").replace("US","$")
            }
        }

        //Anzahl

        let anzahlElement = order.querySelector("div.d-none,d-md-table-cell,text-center:nth-child(1)")
        let anzahl = anzahlElement.textContent
        console.log("Anzahl: " + anzahl)
        /*
        //vendor
        let vendorElement = order.querySelector("div dd a span.PSEUDOLINK")
        console.log(vendorElement.firstChild.textContent)
        */
        //imgUrl
        let imgElement = order.querySelector("div.d-table-cell > div > div.col-2.p-0.p-sm-4 > img")
        let imgUrl = "https://alternate.de" + imgElement.getAttribute("src")
        console.log(imgUrl)




        //build object for later json marshal
        var orderObj = {
            itemName: itemName,
            price: parseInt(price),
            imgUrl: imgUrl,
            purchaseDate: dateElement.textContent,
            anzahl: parseInt(anzahl)
        }
        console.log(JSON.stringify(orderObj))
        //request to api

                postData('http://localhost:8081/order/alternate', orderObj)
          .then(data => {
            console.log(data); // JSON data parsed by `data.json()` call
          });
    }


}, false);