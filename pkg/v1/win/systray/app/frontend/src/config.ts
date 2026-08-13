export const URL = {
    "Domain":"fehmicorp.in",
    "Website": "https://fehmicorp.in",
    "cdn" : "https://cdn.fehmicorp.in/"
}

export const Config = {
    "Name": "Fehmi",
    "Assets":{
        "favicon": `${URL.cdn}fav.ico`,
        "logo":{
            "64": `${URL.cdn}/logo/64.png`,
            "192": `${URL.cdn}/logo/192.png`,
            "512": `${URL.cdn}/logo/512.png`,
        }
    },
    "Company": {
        "LegalName": "Fehmi Corporation",
    },
    "Apps": {
        "Cloud": {
            "Title": "Cloud Connector",
            "Slug": "cloud",
            "Prefix": "/cloud",
            "Tagline":"Managing Cloud makes easy",
            "Description": "High-availability storage and infrastructure management"
        }
    }
} as const;

export type UrlDetails = typeof URL;

export type AppDetails = typeof Config;