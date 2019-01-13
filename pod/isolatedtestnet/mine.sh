#!/bin/bash
podctl -s 127.0.0.1:2000 setgenerate 1 1 &
podctl -s 127.0.0.1:2011 setgenerate 1 1 &
podctl -s 127.0.0.1:2022 setgenerate 1 1 &
podctl -s 127.0.0.1:2033 setgenerate 1 1 &
podctl -s 127.0.0.1:2044 setgenerate 1 1 &
podctl -s 127.0.0.1:2055 setgenerate 1 1 &
podctl -s 127.0.0.1:2066 setgenerate 1 1 &
podctl -s 127.0.0.1:2077 setgenerate 1 1 &
podctl -s 127.0.0.1:2088 setgenerate 1 1 &