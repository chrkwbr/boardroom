const ISO_DATE_TIME_FORMAT = "YYYY-MM-DDTHH:mm:ss.SSSZ";

export const formatDateToIsoDateTime = (date: Date): string => {
  return format(date, ISO_DATE_TIME_FORMAT);
};

const format = (date: Date, format: string): string => {
  const options: Intl.DateTimeFormatOptions = {};
  if (format.includes("YYYY")) {
    options.year = "numeric";
  }
  if (format.includes("MM")) {
    options.month = "2-digit";
  }
  if (format.includes("DD")) {
    options.day = "2-digit";
  }
  if (format.includes("HH")) {
    options.hour = "2-digit";
  }
  if (format.includes("mm")) {
    options.minute = "2-digit";
  }
  if (format.includes("ss")) {
    options.second = "2-digit";
  }
  if (format.includes("Z")) {
    options.timeZoneName = "short";
  }
  return date.toLocaleString("ja-JP", options).replace(/,/g, "-");
};
