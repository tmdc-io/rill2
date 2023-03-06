import {
  FormatterFactory,
  Formatter,
  NumberKind,
  FormatterFactoryOptions,
} from "./humanizer-types";
import { DefaultHumanizer } from "./strategies/default";

export const humanizedFormatterFactory: FormatterFactory = (
  sample: number[],
  options
) => {
  let formatter: Formatter;
  switch (options.strategy) {
    case "default":
      formatter = new DefaultHumanizer(sample, options);
      break;

    default:
      console.warn(
        `Number formatter strategy "${options.strategy}" is not implemented, using default strategy`
      );

      const defaultOptions: FormatterFactoryOptions = {
        strategy: "default",
        padWithInsignificantZeros: true,
        numberKind: options.numberKind || NumberKind.ANY,
        maxDigitsRightSmallNums: 3,
        maxDigitsRightSuffixNums: 2,
      };

      formatter = new DefaultHumanizer(sample, defaultOptions);
      break;
  }

  return formatter;
};
