import streamlit as st
import requests
import urllib.parse

st.title(" Welcome to Web Scraper",)

query = st.text_input("Enter topic of interest")
pages = st.number_input("Number of pages", min_value=1, max_value=5, value=1)

if st.button("Search"):
    if query.strip() == "":
        st.error("Please enter a search topic")
    else:
        encoded_query = urllib.parse.quote(query)
        url = f"http://localhost:8080/search?q={encoded_query}&pages={pages}"

        try:
            response = requests.get(url, timeout=20)

            if response.status_code == 200:
                results = response.json()
                st.success(f"Found {len(results)} results")

                for r in results:
                    st.write(f"### {r['ResultRank']}. {r['ResultTitle']}")
                    st.write(r["ResultURL"])
                    st.write("---")
            else:
                st.error("Backend error: " + response.text)

        except Exception as e:
            st.error(f"Could not connect to backend: {e}")
