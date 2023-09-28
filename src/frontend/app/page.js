import { headers } from 'next/headers'

import Head from 'next/head';
import './globals.css'
import styles from './Style.module.css';


async function getData() {
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
  const res = await fetch(apiUrl, { headers: upstreamHeaders })

  // Pass data to the page via props
  return res.json()
}

export default async function Home() {
  const data = await getData();
  data.sort();
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

        <p className={styles.description}>
          What's in Stock:
        </p>

        <div className={styles.grid}>
          {data.map(item => (
            <a href={`/${item}`} className={styles.card} key={item}>
              <h3>{item}</h3><img src={`${item}.png`} />
            </a>
          ))}
        </div>
      </main>
      <footer>
       Copyright 2023 Beppo's Bottle Shop. All rights reserved.
      </footer>
    </div>
  )
 }
