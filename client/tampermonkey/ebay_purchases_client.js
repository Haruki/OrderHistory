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

(function() {
    'use strict';
    console.log("test")
    let orderList = document.querySelectorAll(".m-order-card")
    console.log(orderList.length)
    for (var order of orderList) {
       //artikelnummer:
       console.log("working on order...")
       let artikelnummer = order.querySelector("[data-listing-id]")
       console.log(artikelnummer.getAttribute("data-listing-id"))
       //itemName:
        let itemElement = order.querySelector("div div div div a.nav-link")
        if (itemElement.length) {
            console.log(itemElement.text)
        }
       //price
        let itemPrice = order.querySelector("div dl dd span.BOLD")
        if (itemPrice.length) {
            console.log(itemPrice.innerHTML.substring(4,).replace(",",""))
        }
        //vendor
        let vendorElement = order.querySelector("div dd a span.PSEUDOLINK")
        console.log(vendorElement.firstChild.textContent)
        //imgUrl
        let imgElement = order.querySelector("div.m-image img[src][alt]")
        console.log(imgElement.getAttribute("src"))
        //orderDate
        let dateElement = order.querySelector("div.ph-col__info-orderDate dd")
        console.log(dateElement.textContent)
    }
})();