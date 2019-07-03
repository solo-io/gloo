import * as React from 'react';
import { isEqual } from 'lodash';

export const useFormValidation = (
  initialState: any,
  validate: any,
  handleSucess: any
) => {
  const [values, setValues] = React.useState(initialState);
  const [errors, setErrors] = React.useState<{ [key: string]: string }>({});
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [isDifferent, setIsDifferent] = React.useState(false);

  React.useEffect(() => {
    if (Object.keys(errors).length === 0 && isSubmitting) {
      handleSucess();
    }
    setIsSubmitting(false);
  }, [errors]);

  // check if values have changed, so we can disable submit for updates
  React.useEffect(() => {
    console.log('initial', initialState);
    console.log('values', values);
    setIsDifferent(!isEqual(values, initialState));
  }, [values]);

  // TODO: type this
  const handleChange = (e: any) => {
    console.log(e.target.value);
    console.log(e.target.name);
    setValues({
      ...values,
      [e.target.name]: e.target.value
    });
  };

  const handleBlur = () => {
    const validationErrors = validate(values);
    setErrors(validationErrors);
  };

  const handleSubmit = (e: React.MouseEvent<any, MouseEvent>) => {
    e.preventDefault();
    setErrors(validate(values));
    setIsSubmitting(true);
  };

  return {
    handleSubmit,
    handleChange,
    handleBlur,
    values,
    errors,
    isSubmitting,
    isDifferent
  };
};
