DontRunKeywords = []
CurrentRunningKeywords = []


def StockConversion(stockStatus) -> int:
    if stockStatus == "IN_STOCK":
        return 1
    elif stockStatus == "OUT_OF_STOCK":
        return 0
    else:
        return -1


def SizeConversion(variantSku: str) -> str:
    variantSizeSplit = variantSku.split(' ')
    if len(variantSizeSplit) > 1:
        return variantSku.split(' ')[1]
    else:
        return variantSku.split('-')[-1]
