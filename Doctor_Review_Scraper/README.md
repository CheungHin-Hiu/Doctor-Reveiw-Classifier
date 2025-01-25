# Doctor Review Scraper

A web scraper built in Go that collects doctor reviews from [goodoctor.com.hk](https://www.goodoctor.com.hk/) using the [colly](https://github.com/gocolly/colly) package. The collected data is intended for machine learning comment classification tasks.

## Prerequisites

- Go 1.x or higher

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/CheungHin-Hiu/Doctor-Reveiw-Classifier.git
   cd doctor_review_scraper
   ```
2. Install dependencies
    ```bash
    go mod download
    ```
## Usage
1. Create a subdirectory called data
    ```bash
    mkdir data\
    ```
2. Run the scraper with optional page range parameters
    ```bash
    go run main.go -start_id={start_page_id} -end_id={end_page_id}
    ```
    - Example: 
        ```bash
        go run main.go -start_id=8037 -end_id=8039
        ```
    - The command will scrape data from page 8037 till 8039
        - https://www.goodoctor.com.hk/doctor/detail/8037
        - https://www.goodoctor.com.hk/doctor/detail/8038
        - https://www.goodoctor.com.hk/doctor/detail/8039
    - Default value:
        - start_id: 6871
        - end_id: 20000

## Output
The scraped data will be saved in the data.csv file in the data directory.