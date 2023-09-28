import { headers } from 'next/headers'

import Head from 'next/head';
import '../globals.css'
import styles from '../Style.module.css';

async function getStock(product) {
  const apiUrl = process.env.API_URL
  const allHeaders = headers()
  var upstreamHeaders = {};
  allHeaders.forEach((value, key) => {
    if (key.startsWith("x-acme")) {
      upstreamHeaders[key] = value
      console.log(`forwarding ${key}`)
    } else {
      console.log(`rejecting ${key}`)
    }
  })
  // Fetch data from external API
  const res = await fetch(`${apiUrl}/${product}`, { headers: upstreamHeaders })
  // Pass data to the page via props
  return res.json()
}

export default async function Page({ params }) {
  const stock = await getStock(params.product)
  return (
      <div>
      <Head>
      <title>Beppo's Bottle Shop</title>
         <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        <h1 className={styles.title}>
          Beppo's Bottle Shop
        </h1>
        <div className={styles.details}>
          <p className={styles.name}>{stock.Product.Name}</p>
          <p className={styles.size}>{stock.Product.Size} {stock.Product.Unit}</p>
          <img src={`${stock.Product.Name}.png`} className={styles.image} />
          <p className={styles.description}>{stock.Product.Description}</p>
          <p className={styles.stock}>{stock.Count} in stock</p>
        </div>
        <p className="homelink">
          <a href="/">Home</a>
        </p>
      </main>
      <footer>
        Copyright 2023 Beppo's Bottle Shop. All rights reserved.
      </footer>
    </div>
  )
}
