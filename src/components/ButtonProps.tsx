import React from 'react';

// ButtonProps interface to define the props for the CustomButton component
interface ButtonProps {
  buttonText: string;
  backgroundColor: string;
  hoverBackgroundColor: string;
  textColor: string;
  additionalClasses?: string;
  onClick: (...args: any[]) => void;
  onClickArgs?: any[];
}

const CustomButton: React.FC<ButtonProps> = ({
  buttonText,
  backgroundColor,
  hoverBackgroundColor,
  textColor,
  additionalClasses = '',
  onClick,
  onClickArgs = [],
}) => {
  const bgClass = `bg-${backgroundColor}`;
  const hoverBgClass = `hover:bg-${hoverBackgroundColor}`;
  const textClass = `text-${textColor}`;
  const classes = `rounded-md py-2 px-3 ${bgClass} ${hoverBgClass} ${textClass} ${additionalClasses}`;

  return (
    <div className="flex justify-center mt-6 mb-3">
      <button className={classes} onClick={() => onClick(...onClickArgs)}>
        {buttonText}
      </button>
    </div>
  );
};

export default CustomButton;